pipeline {
    agent {
        dockerfile {
            dir "src/github.com/linkernetworks/webservice/jenkins"
        }
    }
    options {
        timestamps()
        timeout(time: 1, unit: 'HOURS')
        checkoutToSubdirectory('src/github.com/linkernetworks/web-services')
    }
    stages {
        stage('Prepare') {
            steps {
                withEnv(["GOPATH+GO=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/web-services") {
                        sh "make pre-build"
                    }
                }
            }
        }
        stage('Build') {
            steps {
                withEnv(["GOPATH+GO=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/web-services") {
                        sh "make build"
                    }
                }
            }
        }
        stage('Test') {
            steps {
                withEnv(["GOPATH+GO=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/web-services") {
                        sh "make test 2>&1 | tee >(go-junit-report > report.xml)"
                        junit "report.xml"
                        sh 'gocover-cobertura < coverage.txt > cobertura.xml'
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
}