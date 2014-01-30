var gac = angular.module('getacar', []);

gac.factory('carsFactory', function($http, $q){
  var baseUrl = '/comments';
  return {
    save: function(comment) {
      var deferred = $q.defer();

      $http.post(baseUrl, comment)
      .success(function(data, status, headers, config){
        deferred.resolve(data);
      })
      .error(function(data, status, headers, config){
        deferred.reject(status)
      });

      return deferred.promise;
    },
    query: function() {
      var deferred = $q.defer();

      $http.get("/cars")
      .success(function(data, status, headers, config){
        deferred.resolve(data);
      })
      .error(function(data, status, headers, config){
        deferred.reject(status)
      });

      return deferred.promise;
    },
    checkLogin: function(){
      var deferred = $q.defer();

      $http.get(baseUrl+"/check-login")
      .success(function(data, status, headers, config){
        deferred.resolve(data);
      })
      .error(function(data, status, headers, config){
        deferred.reject(status)
      });

      return deferred.promise;
    }
  };
});


gac.controller('CarsController', function($scope, carsFactory) {
  carsFactory.query().then(function(data){
    $scope.cars = data;
  });

  $scope.addComment = function(){
    $scope.newComment.Url = $attrs.url
    commentsFactory.save($scope.newComment).then(function(){
      $scope.newComment.Avatar = $scope.userSession.Avatar
      $scope.newComment.User = $scope.userSession.Username
      $scope.newComment.Date = new Date();
      $scope.comments.push($scope.newComment)
      $scope.newComment = {}
    });
  }

});

