pipeline {
    agent { label 'light' }
    options {
        timeout(time: 15, unit: 'MINUTES')
        buildDiscarder(logRotator(numToKeepStr: '8'))
        disableConcurrentBuilds()
    }
    stages {
        stage('Build') {
            steps {
                echo 'Building'
               // sh 'go build'
            }
        }
        stage('TestPullRequestDevelop') {
            when {
               allOf {
                    changeRequest()
                }
            }
        	steps {
        		echo 'Running tests'
              //  sh 'make test'
        	}
        }
    }
}
