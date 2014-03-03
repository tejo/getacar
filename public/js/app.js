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
      when('/cars/:lat/:lon/:street', {
        templateUrl: 'cars.html',
        controller: 'CarsController'
      }).
      when('/cars/:lat/:lon/:city/:street', {
        templateUrl: 'cars.html',
        controller: 'CarsController'
      }).
      when('/cars', {
        templateUrl: 'cars.html',
        controller: 'CarsController'
      }).
      when('/map/:carId/:lat/:lon', {
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
});


gac.factory('googleMapsFactory', function() {
  return {
    removeMarkers: function() {
      for( var i = 0; i < this.markers.length; i++ ){
        this.markers[i].setMap(null);
      }
    },
    markers: [],
    marker_limit: 30,
    reference : null,
    init: function(car_array, my_car_id, DataSharingObject) {
      
      /*
      var overlay;
      ColorOverlay.prototype = new google.maps.OverlayView();
      ColorOverlay.prototype.onAdd = function() {
        var div = document.createElement('div');
        div.className = 'map-overlay';
        div.style.width = screen.width;
        div.style.height = screen.height;
        
        this.div_ = div;
        var panes = this.getPanes();
        panes.overlayLayer.appendChild(div);
      }
      ColorOverlay.prototype.draw = function() {

        var div = this.div_;
        div.style.left = '0px';
        div.style.top = '0px';
        div.style.width = ($(window).width()*2) + 'px';
        div.style.height = ($(window).height()*2) + 'px';
      }
      */
    
      var myLatLng = new google.maps.LatLng(DataSharingObject.myLat, DataSharingObject.myLon);
      var map_options = {
        zoom: 16,
        zoomControlOptions: {
          style: google.maps.ZoomControlStyle.LARGE,
          position: google.maps.ControlPosition.RIGHT_CENTER
        },
        panControl: false,
        mapTypeControl: false,
        streetViewControl: false,
        mapTypeId: google.maps.MapTypeId.ROADMAP,
        //center: new google.maps.LatLng(car_array[my_car_id].lat, car_array[my_car_id].lon)
        center: myLatLng
      }
      this.reference = new google.maps.Map(document.getElementById('map-canvas'),map_options);

      var trafficLayer = new google.maps.TrafficLayer();
      trafficLayer.setMap(this.reference);
      
      var myPosition = new google.maps.Marker({
        position: myLatLng,
        map: this.reference,
        //size: new google.maps.Size(20, 32),
        icon: '/images/iamhere.png'
      });
      
      // add colored overlay
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
      //overlay = new ColorOverlay(this.reference);
      
      return this.addCarMarkers(car_array, my_car_id, DataSharingObject);
    },
    addCarMarkers: function(car_array, my_car_id, DataSharingObject) {
      for(var i=0, len=this.marker_limit; i<len; i++) {
        var type = car_array[i].Type;
        var marker = null;
        var image = '/images/car_'+type+'.png';
        if (i == my_car_id) {
          image = '/images/car_car2go.png';
        }
        marker = new google.maps.Marker({
          position: new google.maps.LatLng(car_array[i].lat,car_array[i].lon),
          icon: image,
          animation: google.maps.Animation.DROP,
          map:this.reference
        });
        this.markers.push(marker);
        //this.addListener(marker, DataSharingObject, i);
      }
      return this.markers;
    },
    getMarkerList: function() {
      return this.markers;
    },
    addListener: function(marker, DataSharingObject, car_index) {
      google.maps.event.addListener(marker, 'click', function() {
        //$scope.car.fuel_level = DataSharingObject.cars[car_index].fuel_level;
      });
    }
  }
});


gac.controller('MenuController', function($scope, $location, $routeParams, carsFactory) {
});


gac.controller('HomeController', function($scope, $location, $routeParams, carsFactory, DataSharingObject) {
  
  // Reset
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
  
});


gac.controller('MapController', function($scope, $location, $routeParams, carsFactory, DataSharingObject, googleMapsFactory) {
  
  $scope.lat = $routeParams.lat;
  $scope.lon = $routeParams.lon;
  $scope.getDistance = function(lat, lon) {
    var distance =
      ((calculateDistance($scope.lat, $scope.lon, lat, lon)*1000)
              .toPrecision(4)).toString() + "mt";
    return distance;
  }
  
  if (!DataSharingObject.cars) {
    carsFactory.query($routeParams.lat, $routeParams.lon).then(function(data){
      DataSharingObject.cars = data;
      $scope.cars = DataSharingObject.cars;
      $scope.car = $scope.cars[$routeParams.carId];
      $scope.distance = $routeParams.dist;
      $scope.myLat = $routeParams.lat;
      $scope.myLon = $routeParams.lon;
      $scope.fuel = $scope.car.fuel_level;
      
      var markers = googleMapsFactory.init($scope.cars, $routeParams.carId, DataSharingObject);
      
      for(var i=0, len=markers; i<len.length; i++) {
      
        (function(index){google.maps.event.addListener(markers[index], 'click', function() {
          $scope.$apply(
            function(){$scope.fuel = DataSharingObject.cars[index].fuel_level;}
          );
        });})(i)
        
      }
      

      
      switch($scope.car.Type) {
        case 'car2go':
          DataSharingObject.price = 0.29
        break;
        case 'enjoy':
          DataSharingObject.price = 0.25
        break;
        default:
          0
      }

      // aggiungo attributo prezzo (levare da DataSharingObject?)
      $scope.price = DataSharingObject.price;
      
    });
  } 
  else {
    $scope.cars = DataSharingObject.cars;
    $scope.car = $scope.cars[$routeParams.carId];
    $scope.distance = $routeParams.dist;
    $scope.myLat = $routeParams.lat;
    $scope.myLon = $routeParams.lon;
    $scope.fuel = $scope.car.fuel_level;
    
    var markers = googleMapsFactory.init($scope.cars, $routeParams.carId, DataSharingObject);
    
    for(var i=0, len=markers; i<len.length; i++) {
    
      (function(index){google.maps.event.addListener(markers[index], 'click', function(e) {
        $scope.$apply(
          function(){
            $scope.fuel = DataSharingObject.cars[index].fuel_level;
            var carLat = e.latLng.lat().toFixed(3);
            var carLon = e.latLng.lng().toFixed(3);
            $scope.distance = $scope.getDistance(carLat,carLon);

          }
        );
      });})(i)
      
    }
    
    /*
    
    for(var i=0, len=m; i<len.length; i++) {
      google.maps.event.addListener(m[i], 'click', function(innerKey) {
        $scope.fuel = DataSharingObject.cars[innerKey].fuel_level;
        return function() {
            //$scope.car.fuel_level = DataSharingObject.cars[innerKey].fuel_level;
            $scope.fuel = DataSharingObject.cars[innerKey].fuel_level;

            
        }
      }(i));
    }*/
    
        
    switch($scope.car.Type) {
      case 'car2go':
        DataSharingObject.price = 0.29
        break;
      case 'enjoy':
        DataSharingObject.price = 0.25
        break;
      default:
        0
    }
      
    // aggiungo attributo Price
    $scope.price = DataSharingObject.price;
  }
  
  
});

/*
  DataSharingObject [campi valorizzati]:
  - myLat
  - myLon
  - myCity
  - myStreet
  - price
*/
gac.controller('CarsController', function($scope, $location, $routeParams, carsFactory, DataSharingObject) {
  
  
  // Salvo la posizione (mia/cercata)
  //DataSharingObject.myLat = $routeParams.lat;
  //DataSharingObject.myLon = $routeParams.lon;
  
  $scope.city = $routeParams.city;//DataSharingObject.myCity;
  $scope.street = $routeParams.street;//DataSharingObject.myStreet;
    
  $scope.queryByClosest = function(){
    $scope.lat = $routeParams.lat;
    $scope.lon = $routeParams.lon;
    
    carsFactory.query($scope.lat, $scope.lon).then(function(data){
      $scope.cars = data;
      DataSharingObject.cars = data;
    });
  }();
  $scope.getDistance = function(lat, lon) {
    var distance =
      ((calculateDistance($scope.lat, $scope.lon, lat, lon)*1000)
              .toPrecision(4)).toString() + "mt";
    return distance;
  }
});

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


function ColorOverlay(map) {
  this.map_ = map;
  this.div_ = null;
  this.setMap(map);
}