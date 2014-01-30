var gac = angular.module('getacar', []);

gac.factory('carsFactory', function($http, $q){
  return {
    query: function(lat, lng) {
      var carsUrl = "/cars"
      if(lat != undefined && lng != undefined){
        carsUrl = "/cars/" + lat + "/" + lng
      }
      var deferred = $q.defer();

      $http.get(carsUrl)
      .success(function(data, status, headers, config){
        deferred.resolve(data);
      })
      .error(function(data, status, headers, config){
        deferred.reject(status)
      });

      return deferred.promise;
    },
    geoCode: function(addr) {
      var deferred = $q.defer();

      $http.get("/geocode", {params: {q: addr}})
      .success(function(data, status, headers, config){
        deferred.resolve(data);
      })
      .error(function(data, status, headers, config){
        deferred.reject(status)
      });

      return deferred.promise;
    },
  };
});


gac.controller('CarsController', function($scope, carsFactory) {
  carsFactory.query().then(function(data){
    $scope.cars = data;
  });

  $scope.geoCode = function(){
    carsFactory.geoCode($scope.addrToGeocode).then(function(data){
      console.log(data)
      $scope.queryByClosest(data)
    });
  }

  $scope.queryByClosest = function(geo){
    if(!geo.Success){
      return
    }
    carsFactory.query(geo.Lat, geo.Lng).then(function(data){
      $scope.cars = data;
    });
  }

});

