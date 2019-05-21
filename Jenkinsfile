pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        echo 'Build Image'
        sh 'docker build -t moove/location_service:"build-$BUILD_NUMBER" moove/location_service:latest 482532497705.dkr.ecr.ap-south-1.amazonaws.com/location_service:build-${BUILD_NUMBER} .'
      }
    }
    stage('Push To ECS') {
      steps {
        echo 'Login To ECS Repository'
        sh 'eval $(aws ecr get-login --no-include-email | sed \'s|https://||\')'
        echo 'Push To ECS'
        sh 'docker push 482532497705.dkr.ecr.ap-south-1.amazonaws.com/location_service:build-${BUILD_NUMBER}'
      }
    }
  }
}