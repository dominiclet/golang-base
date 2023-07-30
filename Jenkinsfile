pipeline {
    agent any
    tools {
        go '20'
    }

    stages {
        stage('Build') {
            steps {
                sh 'go build -o output-binary cmd/main.go'
            }
        }
        stage('Stop server') {
        }
        stage('Stage files') {
        }
        stage('Deploy') {
        }
    }
}
