import os

from invoke import task, run
import tempfile


@task
def serve():
    """ run the server """
    run("cd server && GKE_NAMESPACE=local ENVIRONMENT=local go run main.go")


@task
def test(verbose=False):
    """ Runs tests for event-store """
    args = ''
    if verbose:
        args += '-v'
    run("go test {} ./pkg/...".format(args))


@task
def lint():
    """ Run golint on the pkg directory"""
    run("find ./pkg -name '*.go' -and -not -name '*.pb.go' -and -not -name '*_test.go' -exec golint {} \; | sed '/returns unexported type/d'")


@task
def compile_api():
    """ Compiles the Google Cloud Endpoints API files """
    run("docker run --rm -v $PWD:/src vendasta/protoc-python:4.0.1 --proto_path=./vendor/github.com/vendasta/gosdks/pb/event-store/v1/ --include_source_info --include_imports --descriptor_set_out=./vendor/github.com/vendasta/gosdks/pb/event-store/v1/api.descriptor ./vendor/github.com/vendasta/gosdks/pb/event-store/v1/event.proto")

    temp_dir = tempfile.mkdtemp(suffix="api-compiler")
    run("cd {} && git clone https://github.com/googleapis/api-compiler.git".format(temp_dir))

    #for env in ['test', 'demo', 'prod', 'local']:
    for env in ['local']:
        run("cd {temp_dir}/api-compiler && ./run.sh --configs {cur_dir}/endpoints/{env}/event-store.yml --configs {cur_dir}/endpoints/event-store-http.yml --configs {cur_dir}/endpoints/endpoints.yaml --descriptor {cur_dir}/vendor/github.com/vendasta/gosdks/pb/event-store/v1/api.descriptor --json_out {cur_dir}/endpoints/{env}/event-store-service.json".format(
                temp_dir=temp_dir, cur_dir=os.getcwd(), env=env))

@task
def python_protobuf():
    """ Compiles protos into python """
    run("docker run --rm -v $PWD:/src vendasta/protoc-python:4.0.1 --proto_path=./vendor/github.com/vendasta/gosdks/pb/ --python_out=./sdks/python/event-store/_generated/grpc --grpc_python_out=./sdks/python/event-store/_generated/grpc ./vendor/github.com/vendasta/gosdks/pb/event-store/v1/event.proto")
    run("docker run --rm -v $PWD:/src vendasta/protoc-python:4.0.1 --proto_path=./vendor/github.com/vendasta/gosdks/pb/ --python_out=./sdks/python/event-store/_generated/proto/ ./vendor/github.com/vendasta/gosdks/pb/event-store/v1/event.proto")

# gcloud beta service-management
@task
def deploy_endpoints(env="local"):
    """ Deploy cloud endpoints"""
    run("gcloud beta service-management deploy ./vendor/github.com/vendasta/gosdks/pb/event-store/v1/api.descriptor ./endpoints/{}/event-store.yml ./endpoints/event-store-http.yml".format(env))
