---
title: Your short, descriptive title
authors:
- "@robot" # Authors' github accounts here.
reviewers:
- "@robot"
- TBD
approvers:
- "@robot"
- TBD

creation-date: yyyy-mm-dd

---

# Your short, descriptive title

<!--
This is the title of your KEP. Keep it short, simple, and descriptive. A good
title can help communicate what the KEP is and should be considered as part of
any review.
-->

## Summary

With the development of multi-cluster technology, some users begin to want to migrate stateful applications to multi-cluster scenarios, which cannot be covered well by the current version of karmada. Each user needs to implement a `karmada-operator` that does something similar, so we can propose an API to handle the common coordination related work logic, while setting aside a mechanism for the user's `karmada-operator` to complete the specific logic, such as rolling upgrades, scaling, etc.

## Motivation

Provides a new API for stateful applications across clusters, enabling karmada to adapt to stateful applications across clusters.

### Goals

- Defining an API(TBD,karmada/StatefulSet?) enables users to implement specific logic for stateful operator applications across clusters.
-  Propose the implementation ideas for involved components, including the new controller of `cross_cluster_statefulset_controller` in karma -controller-manager.

### Non-Goals

## Proposal

<!--
This is where we get down to the specifics of what the proposal actually is.
This should have enough detail that reviewers can understand exactly what
you're proposing, but should not include things like API designs or
implementation. What is the desired outcome and how do we measure success?.
The "Design Details" section below is for the real
nitty-gritty.
-->

### User Stories (Optional)

#### As a user, I want to deploy a cross cluster stateful app.

I want to let my kubernetes operator can working for cross cluster stateful with karmada and not need to implement other operator for karmada.

#### Story 1

#### Story 2

### Notes/Constraints/Caveats (Optional)

<!--
What are the caveats to the proposal?
What are some important details that didn't come across above?
Go in to as much detail as necessary here.
This might be a good place to talk about core concepts and how they relate.
-->

### Risks and Mitigations

<!--
What are the risks of this proposal, and how do we mitigate? 

How will security be reviewed, and by whom?

How will UX be reviewed, and by whom?

Consider including folks who also work outside the SIG or subproject.
-->

## Design Details

<!--
This section should contain enough information that the specifics of your
change are understandable. This may include API specs (though not always
required) or even code snippets. If there's any ambiguity about HOW your
proposal will be implemented, this is the place to discuss them.
-->

A new API

### Test Plan

<!--
**Note:** *Not required until targeted at a release.*

Consider the following in developing a test plan for this enhancement:
- Will there be e2e and integration tests, in addition to unit tests?
- How will it be tested in isolation vs with other components?

No need to outline all test cases, just the general strategy. Anything
that would count as tricky in the implementation, and anything particularly
challenging to test, should be called out.

-->

## Alternatives

<!--
What other approaches did you consider, and why did you rule them out? These do
not need to be as detailed as the proposal, but should include enough
information to express the idea and why it was not acceptable.
-->

<!--
Note: This is a simplified version of kubernetes enhancement proposal template.
https://github.com/kubernetes/enhancements/tree/3317d4cb548c396a430d1c1ac6625226018adf6a/keps/NNNN-kep-template
-->