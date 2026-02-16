# API Reference

## Packages
- [traefik.extensions.gardener.cloud/v1alpha1](#traefikextensionsgardenercloudv1alpha1)


## traefik.extensions.gardener.cloud/v1alpha1

Package v1alpha1 provides the v1alpha1 version of the external API types.





#### TraefikConfigSpec



TraefikConfigSpec defines the desired state of [TraefikConfig]



_Appears in:_
- [TraefikConfig](#traefikconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image is the Traefik container image to use. |  |  |
| `replicas` _integer_ | Replicas is the number of Traefik replicas to deploy.<br />Defaults to 2 if not specified. |  |  |
| `ingressClass` _string_ | IngressClass is the ingress class name that Traefik will handle.<br />Defaults to "traefik" if not specified.<br />This replaces the deprecated nginx ingress class. |  |  |


