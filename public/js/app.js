var gac = angular.module('getacar', ['ngRoute']);

gac.config(['$routeProvider','$locationProvider',
  function($routeProvider, $locationProvider) {
    // $locationProvider.html5Mode(true);
    $routeProvider.
      when('/', {
        templateUrl: 'home.html',
        controller: 'HomeController'
      }).
      when('/cars/:lat/:lon', {
        templateUrl: 'cars.html',
        controller: 'CarsController'
      }).
      when('/map/:carId', {
        templateUrl: 'map.html',
        controller: 'MapController'
      }).
      otherwise({
        redirectTo: '/'
      });
  }
]);

gac.factory('DataSharingObject', function(){
  return {};
})

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

gac.controller('HomeController', function($scope, $location, $routeParams, carsFactory) {
  $scope.geoCode = function(){
    carsFactory.geoCode($scope.addrToGeocode).then(function(geo){
      console.log(geo)
      if(!geo.Success){
        return
      }
      $location.path("/cars/"+ geo.Lat +"/"+ geo.Lng)
    });
  }
});

gac.controller('MapController', function($scope, $location, $routeParams, carsFactory, DataSharingObject) {
  $scope.car = DataSharingObject.cars[$routeParams.carId];
});

gac.controller('CarsController', function($scope, $location, $routeParams, carsFactory, DataSharingObject) {
  $scope.queryByClosest = function(){
    $scope.lat = $routeParams.lat;
    $scope.lon = $routeParams.lon;
    carsFactory.query($scope.lat, $scope.lon).then(function(data){
      $scope.cars = data;
      DataSharingObject.cars = data;
    });
  }();

  $scope.getDistance = function(lat, lon) {
      return ((calculateDistance($scope.lat, $scope.lon, lat, lon)*1000)
              .toPrecision(4)).toString() + " mt."
    }
});

function calculateDistance(lat1,lon1,lat2,lon2) {
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


