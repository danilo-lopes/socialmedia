pipeline {
  agent any

  stages {
    stage('PR - Teste Unitario') {
      when { branch 'PR-*' }

      steps {
        echo "esse step e executado apenas quando tiver um PR aberto"
        echo env.BRANCH_NAME
      }
    }
  }

  post {
    always {
      node('main'){
        deleteDir()
      }
    }

    success {
      echo "Release Success"
    }

    failure {
      echo "Release Failed"
    }
  }
}
