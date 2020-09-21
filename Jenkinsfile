pipeline {
    agent { label 'light' }
    options {
        timeout(time: 15, unit: 'MINUTES')
        buildDiscarder(logRotator(numToKeepStr: '8'))
        disableConcurrentBuilds()
    }
    stages {
    }
}
