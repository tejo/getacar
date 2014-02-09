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
      when('/map/:carId/:lat/:lon', {
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
    }
  };
});

gac.factory('GoogleMaps', function() {
  return {
    removeMarkers: function() {
      for( var i = 0; i < this.markers.length; i++ ){
        this.markers[i].setMap(null);
      }
    },
    markers: [],
    reference : null,
    init: function(car_array, my_car_id) {
      var map_options = {
        zoom: 16,
        mapTypeId: google.maps.MapTypeId.ROADMAP,
        center: new google.maps.LatLng(car_array[my_car_id].lat, car_array[my_car_id].lon)
      }
      this.reference = new google.maps.Map(document.getElementById('map-canvas'),map_options);
      var trafficLayer = new google.maps.TrafficLayer();
      trafficLayer.setMap(this.reference);
      this.addCarMarkers(car_array, my_car_id);
    },
    addCarMarkers: function(car_array, my_car_id) {
      for(var i=0, len=car_array.length; i<len; i++) {
        var marker = null;
        var image = '/images/ico_other_car.png';
        if (i == my_car_id) {
          image = '/images/ico_my_car.png';
        }
        marker = new google.maps.Marker({
          position: new google.maps.LatLng(car_array[i].lat,car_array[i].lon),
          icon: image,
          animation: google.maps.Animation.DROP,
          map:this.reference
        });
        this.markers.push(marker);
      }
    }
  }
});


gac.controller('MenuController', function($scope, $location, $routeParams, carsFactory) {
  $scope.localizeMe = function() {
    if (Modernizr.geolocation) {
      return navigator.geolocation.getCurrentPosition(
        function(position) {
          $scope.$apply(function() {
            $scope.position = position;
          });
        }, function(error) {
          alert(error);
        }
      );
    } else {
        
    }
  
  
  
  
    console.log(carsFactory.localizeMe().coords);
  }
});


gac.controller('HomeController', function($scope, $location, $routeParams, carsFactory) {
  $scope.geoCode = function(){
    carsFactory.geoCode($scope.addrToGeocode).then(function(geo){
      console.log("home",geo)
      if(!geo.Success){
        return
      }
      $location.path("/cars/"+ geo.Lat +"/"+ geo.Lng)
    });
  }
});

gac.controller('MapController', function($scope, $location, $routeParams, carsFactory, DataSharingObject, GoogleMaps) {
  if (!DataSharingObject.cars) {
    carsFactory.query($routeParams.lat, $routeParams.lon).then(function(data){
      DataSharingObject.cars = data;
      $scope.cars = DataSharingObject.cars;
      $scope.car = $scope.cars[$routeParams.carId];
      GoogleMaps.init($scope.cars, $routeParams.carId);
    });
  }else{
    $scope.cars = DataSharingObject.cars;
    $scope.car = $scope.cars[$routeParams.carId];
    GoogleMaps.init($scope.cars, $routeParams.carId);
  }


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




/*
function loadScript() {
  var script = document.createElement('script');
  script.type = 'text/javascript';
  script.src = 'https://maps.googleapis.com/maps/api/js?key=AIzaSyC1AxcPmZKRn3aNr_1scX7jJS-Ha0uFkaM&sensor=false&' +
      'callback=initialize';
  document.body.appendChild(script);
}
*/

//window.onload = loadScript;

