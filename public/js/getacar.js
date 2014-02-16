// Generated by CoffeeScript 1.6.3
(function() {
  var ButtonsManager, MapManager, OverlayManager, getACar,
    __bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; };

  ButtonsManager = (function() {
    function ButtonsManager() {
      $('form[name=geoCodeForm]').on('.bt-search', 'click', function() {
        return alert('go');
      });
    }

    return ButtonsManager;

  })();

  OverlayManager = (function() {
    function OverlayManager() {
      this.toggleOverlay = __bind(this.toggleOverlay, this);
      this.container = $("div.container");
      this.triggerBttn = $("#trigger-overlay");
      this.overlay = $("div.overlay");
      this.closeBttn = this.overlay.find("button.overlay-close");
      this.transEndEventNames = {
        WebkitTransition: "webkitTransitionEnd",
        MozTransition: "transitionend",
        OTransition: "oTransitionEnd",
        msTransition: "MSTransitionEnd",
        transition: "transitionend"
      };
      this.transEndEventName = this.transEndEventNames[Modernizr.prefixed("transition")];
      this.support = {
        transitions: Modernizr.csstransitions
      };
      this.triggerBttn.on("click", this.toggleOverlay);
      this.closeBttn.on("click", this.toggleOverlay);
    }

    OverlayManager.prototype.toggleOverlay = function() {
      var onEndTransitionFn, _clicked,
        _this = this;
      _clicked = $(this);
      if (this.overlay.hasClass("open")) {
        this.overlay.removeClass("open");
        this.container.removeClass("overlay-open");
        this.overlay.addClass("close");
        onEndTransitionFn = function(ev) {
          if (_this.support.transitions) {
            if (ev.propertyName !== "visibility") {
              return;
            }
            _clicked.off(_this.transEndEventName, onEndTransitionFn);
          }
          _this.overlay.removeClass("close");
        };
        if (this.support.transitions) {
          this.overlay.on(this.transEndEventName, onEndTransitionFn);
        } else {
          onEndTransitionFn();
        }
      } else if (!this.overlay.hasClass("close")) {
        this.overlay.addClass("open");
        this.container.addClass("overlay-open");
      }
    };

    return OverlayManager;

  })();

  MapManager = (function() {
    function MapManager(map_canvas) {
      this.map_canvas = map_canvas;
      this.map_center = new google.maps.LatLng(-34.397, 150.644);
      this.map_zoom = 8;
      google.maps.event.addDomListener(window, 'load', this.init);
    }

    MapManager.prototype.init = function() {
      var mapOptions;
      console.log('load');
      console.log($('#map-canvas').length);
      mapOptions = {
        center: this.map_center,
        zoom: this.map_zoom
      };
      this.map_handle = new google.maps.Map(this.map_canvas, mapOptions);
      return console.log(this.map_handle);
    };

    /*
    loadScript: () ->
      script = document.createElement('script')
      script.type = 'text/javascript'
      script.src = 'https://maps.googleapis.com/maps/api/js?v=3.exp&sensor=false&' +
        'callback=initialize'
      document.body.appendChild(script)
    */


    return MapManager;

  })();

  getACar = {
    start: function() {
      var btsManager, overLayManager;
      btsManager = new ButtonsManager();
      return overLayManager = new OverlayManager();
    }
  };

  $(function() {
    return getACar.start();
  });

}).call(this);
