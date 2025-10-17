pipeline {
    agent any
    environment {
        IMAGE_NAME = 'edy2010/go_htp/go-http-db'
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // Добавьте эти креды в Jenkins
    }
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        stage('Lint Go') {
            steps {
                sh 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'
                sh '$HOME/go/bin/golangci-lint run'
            }
        }
        stage('Lint Dockerfile') {
            steps {
                sh 'docker run --rm -i hadolint/hadolint < Dockerfile'
            }
        }
        stage('Build Docker Image') {
            steps {
                sh "docker build -t ${IMAGE_NAME}:latest ."
            }
        }
        stage('Scan Image for Vulnerabilities') {
            steps {
                sh "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image --exit-code 1 --severity HIGH,CRITICAL ${IMAGE_NAME}:latest"
            }
        }
        stage('Test') {
            steps {
                sh 'docker compose up -d --build'
                script {
                    def tries = 10
                    def success = false
                    for (int i = 0; i < tries; i++) {
                        try {
                            sh 'curl --fail http://localhost:8080/hello'
                            success = true
                            break
                        } catch (err) {
                            sleep 5
                        }
                    }
                    if (!success) {
                        error("App test failed after ${tries} tries")
                    }
                }
                sh 'docker compose down'
            }
        }
        stage('Push Docker Image') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKERHUB_USER', passwordVariable: 'DOCKERHUB_PASS')]) {
                    sh 'echo $DOCKERHUB_PASS | docker login -u $DOCKERHUB_USER --password-stdin'
                    sh "docker push ${IMAGE_NAME}:latest"
                }
            }
        }
        stage('Deploy') {
            steps {
                // Здесь ваш деплой: например, копирование docker-compose.yml на сервер и запуск через ssh
                // Пример:
                // sh 'scp docker-compose.yml user@yourserver:/path/to/deployment/'
                // sh 'ssh user@yourserver docker compose -f /path/to/deployment/docker-compose.yml up -d'
                echo 'Deploy stage (configure according to your deployment strategy)'
            }
        }
    }
    post {
        always {
            sh 'docker system prune -af || true'
        }
    }
}
