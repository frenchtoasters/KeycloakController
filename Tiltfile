# -*- mode: Python -*-

# envsubst_cmd = "./hack/tools/bin/envsubst"
envsubst_cmd = "envsubst"
tools_bin = "./hack/tools/bin"

#Add tools to path
os.putenv('PATH', os.getenv('PATH') + ':' + tools_bin)

update_settings(k8s_upsert_timeout_secs=60)  # on first tilt up, often can take longer than 30 seconds

# set defaults
settings = {
    "allowed_contexts": [
        "kind-keycloak"
    ],
    "deploy_cert_manager": True,
    "preload_images_for_kind": True,
    "kind_cluster_name": "keycloak",
    "cert_manager_version": "v1.1.0",
    "kubernetes_version": "v1.22.3",
}

keys = ["KEYCLOAK_B64ENCODED_CREDENTIALS"]

# global settings
settings.update(read_json(
    "tilt-settings.json",
    default = {},
))

if settings.get("trigger_mode") == "manual":
    trigger_mode(TRIGGER_MODE_MANUAL)

if "allowed_contexts" in settings:
    allow_k8s_contexts(settings.get("allowed_contexts"))

if "default_registry" in settings:
    default_registry(settings.get("default_registry"))

# Users may define their own Tilt customizations in tilt.d. This directory is excluded from git and these files will
# not be checked in to version control.
def include_user_tilt_files():
    user_tiltfiles = listdir("tilt.d")
    for f in user_tiltfiles:
        include(f)

def append_arg_for_container_in_deployment(yaml_stream, name, namespace, contains_image_name, args):
    for item in yaml_stream:
        if item["kind"] == "Deployment" and item.get("metadata").get("name") == name and item.get("metadata").get("namespace") == namespace:
            containers = item.get("spec").get("template").get("spec").get("containers")
            for container in containers:
                if contains_image_name in container.get("image"):
                    container.get("args").extend(args)


def fixup_yaml_empty_arrays(yaml_str):
    yaml_str = yaml_str.replace("conditions: null", "conditions: []")
    return yaml_str.replace("storedVersions: null", "storedVersions: []")


def validate_auth():
    substitutions = settings.get("kustomize_substitutions", {})
    missing = [k for k in keys if k not in substitutions]
    if missing:
        fail("missing kustomize_substitutions keys {} in tilt-settings.json".format(missing))

tilt_helper_dockerfile_header = """
# Tilt image
FROM golang:1.18 as tilt-helper
# Support live reloading with Tilt
RUN wget --output-document /restart.sh --quiet https://raw.githubusercontent.com/windmilleng/rerun-process-wrapper/master/restart.sh  && \
    wget --output-document /start.sh --quiet https://raw.githubusercontent.com/windmilleng/rerun-process-wrapper/master/start.sh && \
    chmod +x /start.sh && chmod +x /restart.sh
"""

tilt_dockerfile_header = """
FROM gcr.io/distroless/base:debug as tilt
WORKDIR /
COPY --from=tilt-helper /start.sh .
COPY --from=tilt-helper /restart.sh .
COPY manager .
"""

def base64_encode(to_encode):
    encode_blob = local("echo '{}' | tr -d '\n' | base64 - | tr -d '\n'".format(to_encode), quiet=True)
    return str(encode_blob)

def base64_encode_file(path_to_encode):
    encode_blob = local("cat {} | tr -d '\n' | base64 - | tr -d '\n'".format(path_to_encode), quiet=True)
    return str(encode_blob)

def read_file_from_path(path_to_read):
    str_blob = local("cat {} | tr -d '\n'".format(path_to_read), quiet=True)
    return str(str_blob)

def base64_decode(to_decode):
    decode_blob = local("echo '{}' | base64 --decode -".format(to_decode), quiet=True)
    return str(decode_blob)

def kustomizesub(folder):
    yaml = local('hack/kustomize-sub.sh {}'.format(folder), quiet=True)
    return yaml

### MRI customizations

def keycloak():
    version = settings.get("keycloak_version")
    keycloak_uri = "https://raw.githubusercontent.com/keycloak/keycloak-quickstarts/latest/kubernetes/keycloak.yaml"
    cmd = "curl -sSL {} | {} | sed 's,LoadBalancer,ClusterIP,g' | kubectl apply -f -".format(keycloak_uri, envsubst_cmd)
    local(cmd, quiet=True)

def keycloakManager():
    # Apply the kustomized yaml for this provider
    substitutions = settings.get("kustomize_substitutions", {})
    os.environ.update(substitutions)

    yaml = str(kustomizesub("./config/default"))

    # Set up a local_resource build of the provider's manager binary.
    local_resource(
        "manager",
        cmd = 'mkdir -p .tiltbuild;CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \'-extldflags "-static"\' -o .tiltbuild/manager cmd/main.go',
        deps = ["api", "cloud", "config", "controllers", "exp", "feature", "pkg", "go.mod", "go.sum", "main.go", "cmd"]
    )

    dockerfile_contents = "\n".join([
        tilt_helper_dockerfile_header,
        tilt_dockerfile_header,
    ])

    entrypoint = ["sh", "/start.sh", "/manager"]
    extra_args = settings.get("extra_args")
    if extra_args:
        entrypoint.extend(extra_args)

    # Set up an image build for the provider. The live update configuration syncs the output from the local_resource
    # build into the container.
    docker_build(
        ref = "keycloak-manager",
        context = "./.tiltbuild/",
        dockerfile_contents = dockerfile_contents,
        target = "tilt",
        entrypoint = entrypoint,
        only = "manager",
        live_update = [
            sync(".tiltbuild/manager", "/manager"),
            run("sh /restart.sh"),
        ],
        ignore = ["templates"]
    )

    k8s_yaml(blob(yaml))

def waitforsystem():
    local("kubectl wait --for=condition=ready --timeout=300s pod --all")


##############################
# Actual work happens here
##############################

validate_auth()

include_user_tilt_files()

load("ext://cert_manager", "deploy_cert_manager")

if settings.get("deploy_cert_manager"):
    deploy_cert_manager()

keycloak()

keycloakManager()

waitforsystem()
