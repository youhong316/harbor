# Changelog

## v1.6.0 (2018-09-11)

- Support manages Helm Charts: From version 1.6.0, Harbor is upgraded to be a composite cloud-native registry, which supports both image management and helm charts management.
- Support LDAP group: User can import an LDAP/AD group to Harbor and assign project roles to it.
- Replicate images with label filter: Use newly added label filter to narrow down the sourcing image list when doing image replication.
- Migrate multiple databases to one unified PostgreSQL database.

## v1.5.0 (2018-05-07)

- Support read-only mode for registry: Admin can set registry to read-only mode before GC. [Details](https://github.com/vmware/harbor/blob/master/docs/user_guide.md#managing-registry-read-only)
- Label support: User can add label to image/repository, and filter images by label on UI/API. [Details](https://github.com/vmware/harbor/blob/master/docs/user_guide.md#managing-labels)
- Show repositories via Cardview.
- Re-work Job service to make it HA ready.

## v1.4.0 (2018-02-07)

- Replication policy rework to support wildcard, scheduled replication.
- Support repository level description.
- Batch operation on projects/repositories/users from UI.
- On board LDAP user when adding member to a project.

## v1.3.0 (2018-01-04)

- Project level policies for blocking the pull of images with vulnerabilities and unknown provenance.
- Remote certificate verification of replication moved to target level.
- Refined all images to improve security.

## v1.2.0 (2017-09-15)

- Authentication and authorization, implementing vCenter Single Sign On across components and role-based access control at the project level. [Read more](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_overview/introduction.html#projects)
- Full integration of the vSphere Integrated Containers Registry and Management Portal user interfaces. [Read more](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_cloud_admin/)
- Image vulnerabilities scanning.

## v1.1.0 (2017-04-18)

- Add in Notary support
- User can update configuration through Harbor UI
- Redesign of Harbor's UI using Clarity
- Some changes to API
- Fix some security issues in token service
- Upgrade base image of nginx for latest openssl version
- Various bug fixes.

## v0.5.0 (2016-12-6)

- Refactory for a new build process
- Easier configuration for HTTPS in prepare script
- Script to collect logs of a Harbor deployment
- User can view the storage usage (default location) of Harbor.
- Add an attribute to disable normal user to create project
- Various bug fixes.

For Harbor virtual appliance:

- Improve the bootstrap process of ova installation.
- Enable HTTPS by default for .ova deployment, users can download the default root cert from UI for docker client or VCH.
- Preload a photon:1.0 image to Harbor for users who have no internet connection.

## v0.4.5 (2016-10-31)

- Virtual appliance of Harbor for vSphere.
- Refactory for new build process.
- Easier configuration for HTTPS in prepare step.
- Updated documents.
- Various bug fixes.

## v0.4.0 (2016-09-23)

- Database schema changed, data migration/upgrade is needed for previous version.
- A project can be deleted when no images and policies are under it.
- Deleted users can be recreated.
- Replication policy can be deleted.
- Enhanced LDAP authentication, allowing multiple uid attributes.
- Pagination in UI.
- Improved authentication for remote image replication.
- Display release version in UI
- Offline installer.
- Various bug fixes.

## v0.3.5 (2016-08-13)

- Vendoring all dependencies and remove go get from dockerfile
- Installer using Docker Hub to download images
- Harbor base images moved to Photon OS (except for official images from third party)
- New Harbor logo
- Various bug fixes

## v0.3.0 (2016-07-15)

- Database schema changed, data migration/upgrade is needed for previous version.
- New UI
- Image replication across multiple registry instances
- Integration with registry v2.4.0 to support image deletion and garbage collection
- Database migration tool
- Bug fixes

## v0.1.1 (2016-04-08)

- Refactored database schema
- Migrate to docker-compose v2 template
- Update token service to support layer mount
- Various bug fixes

## v0.1.0 (2016-03-11)

Initial release, key features include

- Role based access control (RBAC)
- LDAP / AD integration
- Graphical user interface (GUI)
- Auditting and logging
- RESTful API
- Internationalization
