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
      $scope.geo = data
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

  $scope.getDistance = function(lat, lon) {
    if($scope.geo){
      return (($scope.calculateDistance($scope.geo.Lat, $scope.geo.Lng, lat, lon)*1000).toPrecision(4)).toString() + " mt."
    }
  }


  $scope.calculateDistance = function(lat1,lon1,lat2,lon2) {
    var R, a, c, d, dLat, dLon, lon1, lon2, sin_dlat_2, sin_dlon_2;
    R = 6371;
    dLat = (lat2 - lat1) * Math.PI / 180;
    dLon = (lon2 - lon1) * Math.PI / 180;
    lat1 = lat1 * Math.PI / 180;
    lat2 = lat2 * Math.PI / 180;
    sin_dlon_2 = Math.sin(dLon / 2);
    sin_dlat_2 = Math.sin(dLat / 2);
    a = sin_dlat_2 * sin_dlat_2 + sin_dlon_2 * sin_dlon_2 * Math.cos(lat1) * Math.cos(lat2);
    c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return d = R * c;
  }

});


