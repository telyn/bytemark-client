bytemark-client release process
===============================

bytemark-client has a slightly convoluted release process to ensure that it does not get released accidentally or in a broken state, and to ensure that the release has the highest likelihood of being idential across platforms as possible.

bytemark client's version numbers are in MAJOR.MINOR.PATCH form.

Development normally occurs in a feature or bugfix branch. Once the feature or fix is ready for a merge, merge the current develop back into your branch, then create a changelog entry with a version bump - minor for a new feature, patch for a fix. Once done, open a merge request against the develop branch.

When it comes time to release a new version of the client into the public, one of the maintainers will create a release branch for the version you wish to release. release branches have the naming scheme release-MAJOR.MINOR or release-MAJOR.MINOR.PATCH.

Once the release branch builds successfully, they then open a merge request against master. Once that's accepted, one of the maintainers will push a signed version tag for the merge commit, which will cause the release process to occur. 

Should the release process fail and the .gitlab-ci.yml or .gitlab-ci folder need alteration, further changes can be pushed to the release branch, merged into master, and the tag deleted and updated.
