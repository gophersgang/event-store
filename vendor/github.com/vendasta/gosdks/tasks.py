""" Build Tasks """
from invoke import task, run
import os
import glob

@task
def test(verbose=False):
    """ Runs tests for the vStore client """
    args = ''
    if verbose:
        args += '-v'
    run("go test {} ./vstore".format(args))


def protofiles_in_dir(ospath):
    """ Returns the list of protofiles in a path """
    return glob.glob(os.path.join(ospath, "*.proto"))


def generate_protos(ospath):
    """ Recursively generates any proto files inside the given directory """
    files = protofiles_in_dir(ospath)
    if files:
        print("Generating protofiles in directory %s" % ospath)
        cmd = "docker run -v {}:/src vendasta/protoc-go:latest  --go_out=plugins=grpc:. {}".format(ospath, " ".join([os.path.basename(f) for f in files]))
        print("Running command: %s" % cmd)
        run(cmd)
    subdirs = [os.path.join(ospath, name) for name in os.listdir(ospath) if os.path.isdir(os.path.join(ospath, name))]
    return [generate_protos(d) for d in subdirs]


@task
def protobuf():
    """ Generate protobuf files """
    generate_protos(os.path.join(os.getcwd(), "pb"))


@task
def build(verbose=False, branch=None):
    """ Build dependencies, you should only need to run this if you need new or updated protos """
    # reclone vendastaapis protos
    run("rm -rf ./pb")
    branch = branch or "master"
    run("git clone https://github.com/vendasta/vendastaapis --branch %s --single-branch ./pb" % branch)

    # strip vcs metadata
    run('( find ./pb -type d -name ".git" && find ./pb -name ".gitignore" && find ./pb -name ".gitmodules" ) | xargs rm -rf')

    # generate proto files
    protobuf()


@task
def lint(project_dir="./"):
    """ Run golint on the target directory"""
    if not project_dir.startswith("./"):
        project_dir = "./" + project_dir
    run("find %s -not -path '*/vendor/*' -and -name '*.go' -and -not -name '*.pb.go' -and -not -name '*_test.go' -exec golint {} \;" % project_dir)

@task
def mscli_dockerfile(version):
    """ Builds the mscli dockerfile"""
    run("docker build -t vendasta/mscli:{} -f dockerfiles/mscli/Dockerfile .".format(version))
    run("gcloud docker push vendasta/mscli:{}".format(version))
