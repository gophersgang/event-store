
#!groovy
properties ([
      disableConcurrentBuilds()
])

def label = "event-store.${env.BRANCH_NAME}.${env.BUILD_NUMBER}".replace('-', '_').replace('/', '_')
def appCredentials = "${env.GOOGLE_APPLICATION_CREDENTIALS}"
podTemplate(label: label, containers: [
    containerTemplate(
        name: 'jnlp',
        image: 'jenkinsci/jnlp-slave:2.62-alpine',
        args: '${computer.jnlpmac} ${computer.name}',
        workingDir: '/home/jenkins',
        resourceRequestCpu: '.5',
        resourceLimitCpu: '1',
        resourceRequestMemory: '1Gi',
        resourceLimitMemory: '2Gi',
    ),
    containerTemplate(
        name: 'mscli',
        image: 'vendasta/mscli:1.16.0',
        ttyEnabled: true,
        command: 'cat',
        alwaysPullImage: true,
        workingDir: '/home/jenkins',
        resourceRequestCpu: '2',
        resourceLimitCpu: '2',
        resourceRequestMemory: '4Gi',
        resourceLimitMemory: '4Gi',
        envVars: [
            containerEnvVar(key: 'GOOGLE_APPLICATION_CREDENTIALS', value: '/etc/hal9000/hal9000.json'),
            containerEnvVar(key: 'DOCKER_API_VERSION', value: '1.23'),
        ]
    )],
    volumes: [
        secretVolume(mountPath: '/etc/hal9000', secretName: 'hal9000'),
        emptyDirVolume(mountPath: '/home/jenkins'),
        hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
    ]) {
    node(label) {
        def appName = 'event-store-build'
        def imageTag = "gcr.io/repcore-prod/${appName}:${env.BUILD_NUMBER}"

        stage("Checkout") {
            checkout scm
        }

        stage("Build") {
            container('mscli') {
                sh """
                    cp -r /var/run/secrets/kubernetes.io/serviceaccount/ ./serviceaccount
                    ls /etc/hal9000/
                    cat /etc/hal9000/hal9000.json
                    cp /etc/hal9000/hal9000.json .
                    docker build -f ci/Dockerfile -t ${imageTag} .
                """
            }
        }
        stage("Compile") {
            container('mscli') {
                sh("docker run -v /var/run/docker.sock:/var/run/docker.sock --env DOCKER_API_VERSION=1.23 --env JENKINS_HOME=${JENKINS_HOME} --rm --workdir=/go/src/github.com/vendasta/event-store ${imageTag} mscli build ${env.BUILD_NUMBER}")
            }
        }
        stage("Tests") {
            container('mscli') {
                sh("docker run -v /var/run/docker.sock:/var/run/docker.sock --env DOCKER_API_VERSION=1.23 --env JENKINS_HOME=${JENKINS_HOME} --rm --workdir=/go/src/github.com/vendasta/event-store ${imageTag} mscli test ${env.BUILD_NUMBER}")
            }
        }
        stage("Vet") {
            container('mscli') {
                sh("docker run -v /var/run/docker.sock:/var/run/docker.sock --env DOCKER_API_VERSION=1.23 --env JENKINS_HOME=${JENKINS_HOME}  --workdir=/go/src/github.com/vendasta/event-store ${imageTag} mscli vet ${env.BUILD_NUMBER}")
            }
        }
        stage("Lint") {
            container('mscli') {
               sh("docker run -v /var/run/docker.sock:/var/run/docker.sock --env DOCKER_API_VERSION=1.23 --env JENKINS_HOME=${JENKINS_HOME} --workdir=/go/src/github.com/vendasta/event-store ${imageTag} mscli lint ${env.BUILD_NUMBER}")
            }
        }
        switch (env.BRANCH_NAME) {
        case "master":
            stage('Deploy Application to test') {
                container('mscli') {
                    sh("docker run -v /var/run/docker.sock:/var/run/docker.sock -v /var/run/secrets/kubernetes.io/serviceaccount/ --env JENKINS_HOME=${JENKINS_HOME} -e DOCKER_API_VERSION=1.23 -e KUBERNETES_SERVICE_HOST=${env.KUBERNETES_SERVICE_HOST} -e KUBERNETES_SERVICE_PORT=${env.KUBERNETES_SERVICE_PORT} -e GOOGLE_APPLICATION_CREDENTIALS=\"/etc/hal9000/hal9000.json\" --workdir=/go/src/github.com/vendasta/event-store ${imageTag} mscli deploy ${env.BUILD_NUMBER} --env=test")
                }
            }
            stage('Deploy Application to demo') {
                container('mscli') {
                    sh("docker run -v /var/run/docker.sock:/var/run/docker.sock -v /var/run/secrets/kubernetes.io/serviceaccount/ --env JENKINS_HOME=${JENKINS_HOME} -e DOCKER_API_VERSION=1.23 -e KUBERNETES_SERVICE_HOST=${env.KUBERNETES_SERVICE_HOST} -e KUBERNETES_SERVICE_PORT=${env.KUBERNETES_SERVICE_PORT} -e GOOGLE_APPLICATION_CREDENTIALS=\"/etc/hal9000/hal9000.json\" --workdir=/go/src/github.com/vendasta/event-store ${imageTag} mscli deploy ${env.BUILD_NUMBER} --env=demo")
                }
            }
            stage("Deploy to production?") {
                timeout(5) {
                    input 'Ready to deploy to production?'
                }
            }

            stage("Deploy Application to production"){
                container('mscli') {
                    sh("docker run -v /var/run/docker.sock:/var/run/docker.sock -v /var/run/secrets/kubernetes.io/serviceaccount/ --env JENKINS_HOME=${JENKINS_HOME} -e DOCKER_API_VERSION=1.23 -e KUBERNETES_SERVICE_HOST=${env.KUBERNETES_SERVICE_HOST} -e KUBERNETES_SERVICE_PORT=${env.KUBERNETES_SERVICE_PORT} -e GOOGLE_APPLICATION_CREDENTIALS=\"/etc/hal9000/hal9000.json\" --workdir=/go/src/github.com/vendasta/event-store ${imageTag} mscli deploy ${env.BUILD_NUMBER} --env=prod")
                }
            }
            break
        }
    }
}
