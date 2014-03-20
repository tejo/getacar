// @codekit-prepend "angular.min.js"
// @codekit-prepend "angular-route.min.js"
// @codekit-append "getacar.min.js"

var gac = angular.module('getacar', ['ngRoute']);

gac.config(['$routeProvider','$locationProvider',
  function($routeProvider, $locationProvider) {
    // $locationProvider.html5Mode(true);
    $routeProvider.
      when('/', {
        templateUrl: 'home.html',
        controller: 'HomeController'
      }).
      when('/cars/:lat?/:lon?/:city?/:street?', {
        templateUrl: 'cars.html',
        controller: 'CarsController'
      }).
      when('/map/:carId/:lat/:lon/', {
        templateUrl: 'map.html',
        controller: 'MapController'
      }).
      when('/map/:carId/:lat/:lon/:dist', {
        templateUrl: 'map.html',
        controller: 'MapController'
      });
  }
]);



/*
  @prop myLat
  @prop myLon
  @prop myCity
  @prop myStreet
*/
gac.factory('DataSharingObject', function(){
  return {
    reset : function() {
      this.myLat = '';
      this.myLon = ''
      this.myCity = ''
      this.myStreet = ''
      this.cars = null
    }
  };
})


gac.factory('carsFactory', ['$http', '$q', function($http, $q){
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
    getLocation : function() {
      var deferred;
      deferred = $q.defer();
      if (navigator && navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(function(position) {
          return deferred.resolve(position.coords);
        }, function(error) {
          return deferred.reject("Unable to get your location");
        });
      } else {
        deferred.reject("Your browser cannot access to your position");
      }
      return deferred.promise;
    },
    getAddress : function() {
      var deferred = $q.defer();
      this.geocoder || (this.geocoder = new google.maps.Geocoder());
  
      this.getLocation().then((function(_this) {
        return function(coords) {
          var latlng;
          latlng = new google.maps.LatLng(coords.latitude, coords.longitude);
          return _this.geocoder.geocode({
            latLng: latlng
          }, function(results, status) {
              if (status === google.maps.GeocoderStatus.OK) {
                return deferred.resolve(_this.extractAddress(results, coords.latitude, coords.longitude));
              } else {
                return deferred.reject("cannot geocode status: " + status);
              }
            }, function() {
              return deferred.reject("cannot geocode");
            });
        };
      })(this));
      return deferred.promise;
    },
    extractAddress : function(addresses, lat, lng) {
      var address, component, result, _i, _j, _len, _len1, _ref;
      result = {};
      for (_i = 0, _len = addresses.length; _i < _len; _i++) {
        address = addresses[_i];
        result.fullAddress || (result.fullAddress = address.formatted_address);
        result.coord || (result.coord = [address.geometry.location.ob, address.geometry.location.pb]);
        _ref = address.address_components;
        for (_j = 0, _len1 = _ref.length; _j < _len1; _j++) {
          component = _ref[_j];
          if (component.types[0] === "route") {
            result.street || (result.street = component.long_name);
          }
          if (component.types[0] === "locality") {
            result.city || (result.city = component.long_name);
          }
          if (component.types[0] === "postal_code") {
            result.zip || (result.zip = component.long_name);
          }
          if (component.types[0] === "country") {
            result.country || (result.country = component.long_name);
          }
        }
      }
      result.lat = lat;
      result.lng = lng;
      return result;
    }
    
  };
}]);


gac.factory('googleMapsFactory', function() {
  return {
    markers: [],
    marker_limit: 30,
    reference : null,
    max_cars : 3,
    init: function(car_array, my_car_id, DataSharingObject) {
    
      // Reset markers array
      if (this.markers.length > 0) {
        this.removeMarkers();
      }
      
      var myLatLng = new google.maps.LatLng(DataSharingObject.myLat, DataSharingObject.myLon);
      var styles = [
                      {"featureType":"water","elementType":"geometry","stylers":[{"color":"#4ca6a2"}]},
                      {"featureType":"landscape","elementType":"geometry","stylers":[{"color":"#59b7a3"}]},
                      {"featureType":"road","elementType":"geometry","stylers":[{"color":"#4b9b8c"},{"lightness":0}]},
                      {"featureType":"poi","elementType":"geometry","stylers":[{"color":"#4b9b8c"}]},
                      {"featureType":"transit","elementType":"geometry","stylers":[{"color":"#4b9b8c"}]},
                      {"elementType":"labels.text.stroke","stylers":[{"visibility":"off"},{"color":"#0b3930"},{"weight":2},{"gamma":0.84}]},
                      {"elementType":"labels.text.fill","stylers":[{"color":"#083a30"}]},
                      {"featureType":"administrative","elementType":"geometry","stylers":[{"weight":0.6},{"color":"#000000"}]},
                      {"elementType":"labels.icon","stylers":[{"visibility":"off"}]},
                      {"featureType":"poi.park","elementType":"geometry","stylers":[{"color":"#4b9b8c"}]}
                    ]

      var map_options = {
        mapTypeControlOptions: {
          mapTypeIds: ['Styled']
        },
        zoom: 16,
        disableDefaultUI: true,
        zoomControl: true,
        zoomControlOptions: {
          style: google.maps.ZoomControlStyle.SMALL,
          position: google.maps.ControlPosition.RIGHT_CENTER
        },
        panControl: false,
        mapTypeControl: false,
        streetViewControl: false,
        mapTypeId: 'Styled',
        center: myLatLng
      }

      this.reference = new google.maps.Map(document.getElementById('map-canvas'),map_options);
      var styledMapType = new google.maps.StyledMapType(styles, { name: 'Styled' });  
      this.reference.mapTypes.set('Styled', styledMapType); 

      //var trafficLayer = new google.maps.TrafficLayer();
      //trafficLayer.setMap(this.reference);
      var confini = new google.maps.LatLngBounds();
      
      var myPosition = new google.maps.Marker({
        position: myLatLng,
        map: this.reference,
        icon: '/images/iamhere.png'
        /*size: new google.maps.Size(20, 32),*/
        /*optimized: false,
        icon: '/images/iamhere@x2.gif'*/
      });
      
      // Start extendi map to hold my position
      confini.extend(myPosition.position);
      this.reference.fitBounds(confini);
      
      
      // add colored overlay
      /*
      var bounds = new google.maps.LatLngBounds(
        new google.maps.LatLng(-84.999999, -179.999999),
        new google.maps.LatLng(84.999999, 179.999999)
      );
      rect = new google.maps.Rectangle({
        bounds: bounds,
        fillColor: "#16a085",
        fillOpacity: 0.7,
        strokeWeight: 0,
        map: this.reference
      });
      */
      
      return this.addCarMarkers(car_array, my_car_id, DataSharingObject, confini);
    },
    addCarMarkers: function(car_array, my_car_id, DataSharingObject, confini) {
      for(var i=0, len=this.marker_limit; i<len; i++) {
        var type = car_array[i].Type;
        var marker = null;
        var image = '/images/car_'+type+'.png';
        if (i == my_car_id) {
          image = '/images/ico_selected_car.png';
        }
        marker = new google.maps.Marker({
          position: new google.maps.LatLng(car_array[i].lat,car_array[i].lon),
          icon: image,
          animation: google.maps.Animation.DROP,
          map:this.reference
        });
        this.markers.push(marker);
        if(i<this.max_cars+1) {
          var position = marker.position;
          confini.extend(position);
        }
      }
      // ristringo lo zoom a max_cars
      this.reference.fitBounds(confini);
      
      return this.markers;
    },
    removeMarkers: function() {
      this.markers.length = 0;
    },
    getMarkerList: function() {
      return this.markers;
    }
  }
});


gac.controller('MenuController',['$scope', '$location', '$routeParams', 'carsFactory', function($scope, $location, $routeParams, carsFactory) {
}]);


gac.controller('HomeController', ['$scope', '$location', '$routeParams', 'carsFactory', 'DataSharingObject', function($scope, $location, $routeParams, carsFactory, DataSharingObject) {
  
  // Reset dati condivisi
  DataSharingObject.reset();
  
  /* Auto Localization */
  $scope.localizeMe = function() {
    
    $('.bt-getme').addClass('spin');
    
    carsFactory.getAddress().then(function(address) {
      $scope.address = address;
      DataSharingObject.myCity = address.city;
      DataSharingObject.myStreet = address.street;
      DataSharingObject.myLat = address.lat;
      DataSharingObject.myLon = address.lng;

      $location.path("/cars/"+ address.lat +"/"+ address.lng +"/"+ address.city +"/"+ address.street);
      $('#switcher').addClass('open');
    });
  }
  
  /* Open / Close / Geocode */
  $scope.searchBTMngr = function(e) {
    var _bt = $('form[name=geoCodeForm] .bt-search');
    var _textbox = _bt.siblings('.address-box');

    if( _bt.hasClass('active') ) {

      if( _textbox.val() != "" ) {
        // geocode
        $scope.geoCode();
      }else {
        // close textbox
        _bt.toggleClass('active');
        _textbox.toggleClass('active');
      }
    }else {
      _bt.toggleClass('active');
      _textbox.toggleClass('active');
    }
  }
  
  /* Manual Localization */
  $scope.geoCode = function(){
    carsFactory.geoCode($scope.addrToGeocode).then(function(geo){
      if(!geo.Success){
        return
      }

      DataSharingObject.myStreet = $scope.addrToGeocode;
      DataSharingObject.myLat = geo.Lat;
      DataSharingObject.myLon = geo.Lng;
      
      $location.path("/cars/"+ geo.Lat +"/"+ geo.Lng +"/"+ $scope.addrToGeocode);
      $('#switcher').addClass('open');
    });
  }
  
}]);


gac.directive('featurevalue', ['$timeout',function ($timeout) {
  return {
    restrict: 'C',
    link: function (scope, element, attrs) {
      scope.$watch('fuel', function(newVal, oldVal){
        if (newVal !== oldVal) {
          if (element.hasClass('fuel')) {
            element.addClass("updated");
            $timeout(function(){element.removeClass("updated");}, 1000)
          }
        }
      });
      scope.$watch('distance', function(newVal, oldVal){
        if (newVal !== oldVal) {
          if (element.hasClass('distance')) {
            element.addClass("updated");
            $timeout(function(){element.removeClass("updated");}, 1000)
          }
        }
      });
      scope.$watch('price', function(newVal, oldVal){
        if (newVal !== oldVal) {
          if (element.hasClass('price')) {
            element.addClass("updated");
            $timeout(function(){element.removeClass("updated");}, 1000)
          }
        }
      });            
      /*
scope.featurevalues = ['fuel', 'distance', 'price'];
      scope.$watchCollection('[fuel,distance,price]', function(newVals, oldVals){
        if (newVals !== oldVals) {
          element.addClass("updated");
          $timeout(function(){element.removeClass("updated");}, 1000)
        }
      });
*/ 
/*
      scope.$watchCollection('[fuel,distance,price]', function(newVal, oldVal){
        if (newVal !== oldVal) {
          element.addClass("updated");
          $timeout(function(){element.removeClass("updated");}, 1000)
        }
      })
*/
    }
  };
}]);


gac.controller('MapController', ['$scope', '$location', '$routeParams', 'carsFactory', 'DataSharingObject', 'googleMapsFactory', function($scope, $location, $routeParams, carsFactory, DataSharingObject, googleMapsFactory) {
  
  $scope.lat = $routeParams.lat;
  $scope.lon = $routeParams.lon;
  $scope.getDistance = function(lat, lon) {
    var unit = ' m';
    var distance =
      ((calculateDistance($scope.lat, $scope.lon, lat, lon)*1000)
              .toPrecision(4)).toString();
    if (distance >= 1000) {
      unit = ' km'
      var rounding = preciseRound(distance/1000, 2);
      return rounding + unit;
    }
    return distance + unit;
  }
  
  $scope.city = DataSharingObject.city;
  $scope.street = DataSharingObject.street;

  if (!DataSharingObject.cars) {
    $location.path( "/" );
  } 
  else {
    $scope.cars = DataSharingObject.cars;
    $scope.car = $scope.cars[$routeParams.carId];
    $scope.distance = $routeParams.dist;
    $scope.myLat = $routeParams.lat;
    $scope.myLon = $routeParams.lon;
    $scope.fuel = $scope.car.fuel_level;
    $scope.price = $scope.car.Price;
    
    var markers = googleMapsFactory.init($scope.cars, $routeParams.carId, DataSharingObject);

    for(var i=0, len=markers; i<len.length; i++) {
    
      (function(index){google.maps.event.addListener(markers[index], 'click', function(e) {
        
        $scope.$apply(
          function(){
            
            for (var i=0; i<markers.length; i++) {
              markers[i].setIcon('/images/car_'+DataSharingObject.cars[i].Type+'.png');
            }
            
            markers[index].setIcon('/images/ico_selected_car.png');
            if ($scope.fuel != DataSharingObject.cars[index].fuel_level)
              $scope.fuel = DataSharingObject.cars[index].fuel_level;

            if ($scope.price != DataSharingObject.cars[index].Price)
              $scope.price = DataSharingObject.cars[index].Price;

            var carLat = e.latLng.lat().toFixed(5);
            var carLon = e.latLng.lng().toFixed(5);
            $scope.distance = $scope.getDistance(carLat,carLon);

          }
        );
      });})(i)
      
    }
  }
  
  
}]);

/*
  DataSharingObject [campi valorizzati]:
  - myLat
  - myLon
  - myCity
  - myStreet
  - price
*/
gac.controller('CarsController', ['$scope', '$location', '$routeParams', 'carsFactory', 'DataSharingObject', function($scope, $location, $routeParams, carsFactory, DataSharingObject) {
  
  
  // Salvo la posizione (mia/cercata)
  //DataSharingObject.myLat = $routeParams.lat;
  //DataSharingObject.myLon = $routeParams.lon;
  

  if ($routeParams.city != undefined)
    $scope.city = $routeParams.city;
  else
    $scope.city = "";
  $scope.street = $routeParams.street;
  DataSharingObject.city = $scope.city;//$routeParams.city;
  DataSharingObject.street = $routeParams.street;
    
  $scope.queryByClosest = function(){
    $scope.lat = $routeParams.lat;
    $scope.lon = $routeParams.lon;
    
    carsFactory.query($scope.lat, $scope.lon).then(function(data){
      $scope.cars = data;
      DataSharingObject.cars = data;
    });
  }();
  $scope.getDistance = function(lat, lon) {
    var unit = ' m';
    var distance =
      ((calculateDistance($scope.lat, $scope.lon, lat, lon)*1000)
              .toPrecision(4)).toString();
    if (distance >= 1000) {
      unit = ' km'
      //return (distance/1000) + unit;
      var rounding = preciseRound(distance/1000, 2);
      return rounding + unit;
    }
    return distance + unit;
  }
}]);

/*
  lat1: myLat
  lon1: myLon
  lat2: selected car lat
  lon2: selected car lon
*/
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
  Round distance to two decimals if in kms
*/
function preciseRound(num, decimals) {
  var sign = num >= 0 ? 1 : -1;
  return (Math.round((num*Math.pow(10,decimals))+(sign*0.001))/Math.pow(10,decimals)).toFixed(decimals);
}