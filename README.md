# Valid_Controller
This is a Validating Admission Controller framed for a project (`Digital_Ocean_Cluster_Operator`) that I had hosted before. This means any request made to the Kube-API server with respect to the custom resource `Digital_Ocean_Cluster` will be intercepted by this validation admission component.

### Use Case
1) As the description mentions this is a controller framed for an operator named `Digital_Ocean_Cluster_Operator` that I had built previously . The operator watches for creation and deletion of a custom resource named `Digital_Ocean_Cluster` and simultaneously creates/deletes a K8s cluster of Digital Ocean.
2) Now what could go wrong here is that , when a user might deploy the custom resource , he/she may specify the a version for the Cluster which might not be supported by Digital Ocean . Hence the `Digital_Ocean_Cluster` resource would be deployed on locally but creation of the K8s Cluster would fail on Digital Ocean.
3) Hence it's better to always validate before hand whether the version is comptaible on Digital Ocean or not and maybe return a warning/error message back to the user.
