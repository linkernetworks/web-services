pipeline {
    agent {
        kubernetes {
            cloud 'kubernetes'
            label 'webservice-pod' // must have
            defaultContainer 'golang'
            yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: golang
    image: golang
    command:
    - cat
    tty: true
"""
        }
    }
    options {
        timestamps()
        timeout(time: 1, unit: 'HOURS')
    }
    stages {
        stage('Prepare') {
            steps {
                container('golang') {
                    sh "go get -u github.com/kardianos/govendor"
                    sh "go get -u github.com/jstemmer/go-junit-report"
                    sh "go get -u github.com/t-yuki/gocover-cobertura"
                    sh """
                        mkdir -p /go/src/github.com/linkernetworks
                        ln -s `pwd` /go/src/github.com/linkernetworks/webservice
                        cd /go/src/github.com/linkernetworks/webservice
                        make pre-build
                    """
                }
            }
        }
        stage('Build') {
            steps {
                container('golang') {
                    sh """
                        cd /go/src/github.com/linkernetworks/webservice
                        make build
                    """
                }
            }
        }
        stage('Test') {
            steps {
                container('golang') {

                    sh """
                        cd /go/src/github.com/linkernetworks/webservice
                        make test 2>&1 | tee >(go-junit-report > report.xml)
                    """
                    junit "report.xml"

                    sh """
                        cd /go/src/github.com/linkernetworks/webservice
                        gocover-cobertura < coverage.txt > cobertura.xml
                    """

                    cobertura coberturaReportFile: "cobertura.xml", failNoReports: true, failUnstable: true
                    publishHTML (target: [
                        allowMissing: true,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: './',
                        reportFiles: 'coverage.html',
                        reportName: "GO cover report",
                        reportTitles: "GO cover report",
                        includes: "coverage.html"
                    ])
                }
            }
        }
    }
}