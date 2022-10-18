package proxy

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	"github.com/karmada-io/karmada/pkg/search/proxy/store"
	"github.com/karmada-io/karmada/pkg/util/lifted"
)

var (
	podGVR     = corev1.SchemeGroupVersion.WithResource("pods")
	nodeGVR    = corev1.SchemeGroupVersion.WithResource("nodes")
	secretGVR  = corev1.SchemeGroupVersion.WithResource("secret")
	clusterGVR = clusterv1alpha1.SchemeGroupVersion.WithResource("cluster")
)

func TestCacheProxy_connect(t *testing.T) {
	type args struct {
		url string
	}
	type want struct {
		namespace   string
		name        string
		gvr         schema.GroupVersionResource
		getOptions  *metav1.GetOptions
		listOptions *metainternalversion.ListOptions
	}

	var actual want
	p := &cacheProxy{
		store: &restReaderFuncs{
			GetFunc: func(ctx context.Context, gvr schema.GroupVersionResource, name string, options *metav1.GetOptions) (runtime.Object, error) {
				actual = want{}
				actual.namespace = request.NamespaceValue(ctx)
				actual.name = name
				actual.gvr = gvr
				actual.getOptions = options
				return nil, nil
			},
			ListFunc: func(ctx context.Context, gvr schema.GroupVersionResource, options *metainternalversion.ListOptions) (runtime.Object, error) {
				actual = want{}
				actual.namespace = request.NamespaceValue(ctx)
				actual.gvr = gvr
				actual.listOptions = options
				return nil, nil
			},
			WatchFunc: func(ctx context.Context, gvr schema.GroupVersionResource, options *metainternalversion.ListOptions) (watch.Interface, error) {
				actual = want{}
				actual.namespace = request.NamespaceValue(ctx)
				actual.gvr = gvr
				actual.listOptions = options
				w := newEmptyWatch()
				// avoid block in ServeHTTP
				time.AfterFunc(time.Millisecond*10, func() {
					w.Stop()
				})
				return w, nil
			},
		},
		restMapper: restMapper,
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "get node",
			args: args{
				url: "/api/v1/nodes/foo",
			},
			want: want{
				name:       "foo",
				gvr:        nodeGVR,
				getOptions: &metav1.GetOptions{},
			},
		},
		{
			name: "get pod",
			args: args{
				url: "/api/v1/namespaces/default/pods/foo",
			},
			want: want{
				namespace:  "default",
				name:       "foo",
				gvr:        podGVR,
				getOptions: &metav1.GetOptions{},
			},
		},
		{
			name: "get pod with options",
			args: args{
				url: "/api/v1/namespaces/default/pods/foo?resourceVersion=1000",
			},
			want: want{
				namespace:  "default",
				name:       "foo",
				gvr:        podGVR,
				getOptions: &metav1.GetOptions{ResourceVersion: "1000"},
			},
		},
		{
			name: "list nodes",
			args: args{
				url: "/api/v1/nodes",
			},
			want: want{
				gvr:         nodeGVR,
				listOptions: &metainternalversion.ListOptions{},
			},
		},
		{
			name: "list pod",
			args: args{
				url: "/api/v1/namespaces/default/pods",
			},
			want: want{
				namespace:   "default",
				gvr:         podGVR,
				listOptions: &metainternalversion.ListOptions{},
			},
		},
		{
			name: "list pod with options",
			args: args{
				url: "/api/v1/namespaces/default/pods?fieldSelector=metadata.name%3Dbar&labelSelector=app%3Dfoo&limit=500&resourceVersion=1000&container=bar",
			},
			want: want{
				namespace: "default",
				gvr:       podGVR,
				listOptions: &metainternalversion.ListOptions{
					LabelSelector:   asLabelSelector("app=foo"),
					FieldSelector:   fields.OneTermEqualSelector("metadata.name", "bar"),
					ResourceVersion: "1000",
					Limit:           500,
				},
			},
		},
		{
			name: "watch node",
			args: args{
				url: "/api/v1/nodes?watch=true",
			},
			want: want{
				gvr: nodeGVR,
				listOptions: &metainternalversion.ListOptions{
					LabelSelector: labels.NewSelector(),
					FieldSelector: fields.Everything(),
					Watch:         true,
				},
			},
		},
		{
			name: "watch pod",
			args: args{
				url: "/api/v1/namespaces/default/pods?watch=true",
			},
			want: want{
				namespace: "default",
				gvr:       podGVR,
				listOptions: &metainternalversion.ListOptions{
					LabelSelector: labels.NewSelector(),
					FieldSelector: fields.Everything(),
					Watch:         true,
				},
			},
		},
		{
			name: "watch pod with options",
			args: args{
				url: "/api/v1/namespaces/default/pods?watch=true&fieldSelector=metadata.name%3Dbar&labelSelector=app%3Dfoo&limit=500&resourceVersion=1000&container=bar",
			},
			want: want{
				namespace: "default",
				gvr:       podGVR,
				listOptions: &metainternalversion.ListOptions{
					LabelSelector:   asLabelSelector("app=foo"),
					FieldSelector:   fields.OneTermEqualSelector("metadata.name", "bar"),
					ResourceVersion: "1000",
					Limit:           500,
					Watch:           true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset before each test
			actual = want{}

			req, err := http.NewRequest("GET", tt.args.url, nil)
			if err != nil {
				t.Error(err)
				return
			}
			requestInfo := lifted.NewRequestInfo(req)
			req = req.WithContext(request.WithRequestInfo(req.Context(), requestInfo))
			if requestInfo.Namespace != "" {
				req = req.WithContext(request.WithNamespace(req.Context(), requestInfo.Namespace))
			}

			h, err := p.connect(req.Context(), podGVR, "", nil)
			if err != nil {
				t.Error(err)
				return
			}
			h.ServeHTTP(&emptyResponseWriter{}, req)
			if tt.want.namespace != actual.namespace {
				t.Errorf("want namespace %v, get %v", tt.want.namespace, actual.namespace)
			}
			if tt.want.name != actual.name {
				t.Errorf("want name %v, get %v", tt.want.name, actual.name)
			}
			if tt.want.gvr != actual.gvr {
				t.Errorf("want gvr %v, get %v", tt.want.gvr, actual.gvr)
			}
			if !reflect.DeepEqual(tt.want.getOptions, actual.getOptions) {
				t.Errorf("getOptions diff: %v", cmp.Diff(tt.want.getOptions, actual.getOptions))
			}
			if !reflect.DeepEqual(tt.want.listOptions, actual.listOptions) {
				t.Errorf("listOptions diff: %v", diff.ObjectGoPrintSideBySide(tt.want.listOptions, actual.listOptions))
			}
		})
	}
}

type restReaderFuncs struct {
	GetFunc   func(ctx context.Context, gvr schema.GroupVersionResource, name string, options *metav1.GetOptions) (runtime.Object, error)
	ListFunc  func(ctx context.Context, gvr schema.GroupVersionResource, options *metainternalversion.ListOptions) (runtime.Object, error)
	WatchFunc func(ctx context.Context, gvr schema.GroupVersionResource, options *metainternalversion.ListOptions) (watch.Interface, error)
}

var _ store.RESTReader = &restReaderFuncs{}

func (r *restReaderFuncs) Get(ctx context.Context, gvr schema.GroupVersionResource, name string, options *metav1.GetOptions) (runtime.Object, error) {
	if r.GetFunc == nil {
		panic("implement me")
	}
	return r.GetFunc(ctx, gvr, name, options)
}

func (r *restReaderFuncs) List(ctx context.Context, gvr schema.GroupVersionResource, options *metainternalversion.ListOptions) (runtime.Object, error) {
	if r.GetFunc == nil {
		panic("implement me")
	}
	return r.ListFunc(ctx, gvr, options)
}

func (r *restReaderFuncs) Watch(ctx context.Context, gvr schema.GroupVersionResource, options *metainternalversion.ListOptions) (watch.Interface, error) {
	if r.GetFunc == nil {
		panic("implement me")
	}
	return r.WatchFunc(ctx, gvr, options)
}

type emptyResponseWriter struct{}

var _ http.ResponseWriter = &emptyResponseWriter{}
var _ http.Flusher = &emptyResponseWriter{}

func (n *emptyResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (n *emptyResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (n *emptyResponseWriter) WriteHeader(int) {
}

func (n *emptyResponseWriter) Flush() {
}

type emptyWatch struct {
	ch       chan watch.Event
	isClosed bool
	lock     sync.Mutex
}

func newEmptyWatch() watch.Interface {
	w := &emptyWatch{
		ch: make(chan watch.Event),
	}

	return w
}

func (e *emptyWatch) Stop() {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.isClosed {
		return
	}
	e.isClosed = true
	close(e.ch)
}

func (e *emptyWatch) ResultChan() <-chan watch.Event {
	return e.ch
}

func asLabelSelector(s string) labels.Selector {
	selector, err := labels.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("Fail to parse %s to labels: %v", s, err))
	}
	return selector
}
