// Generated by CoffeeScript 1.6.3
(function() {
  var OverlayManager, dom, getACar,
    __bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; };

  dom = {
    hasClass: function(elem, className) {
      return new RegExp(' ' + className + ' ').test(' ' + elem.className + ' ');
    },
    addClass: function(elem, className) {
      if (!dom.hasClass(elem, className)) {
        return elem.className += ' ' + className;
      }
    },
    removeClass: function(elem, className) {
      var newClass;
      newClass = ' ' + elem.className.replace(/[\t\r\n]/g, ' ') + ' ';
      if (dom.hasClass(elem, className)) {
        while (newClass.indexOf(' ' + className + ' ') >= 0) {
          newClass = newClass.replace(' ' + className + ' ', ' ');
        }
        return elem.className = newClass.replace(/^\s+|\s+$/g, '');
      }
    }
  };

  getACar = {
    start: function() {
      var overLayManager;
      return overLayManager = new OverlayManager();
    }
  };

  OverlayManager = (function() {
    function OverlayManager() {
      this.toggleOverlay = __bind(this.toggleOverlay, this);
      this.container = document.querySelector('div.container');
      this.triggerBttn = document.getElementById('trigger-overlay');
      this.overlay = document.querySelector('div.overlay');
      this.closeBttn = this.overlay.querySelector('button.overlay-close');
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
      this.triggerBttn.addEventListener('click', this.toggleOverlay);
      this.closeBttn.addEventListener('click', this.toggleOverlay);
    }

    OverlayManager.prototype.toggleOverlay = function(e) {
      var onEndTransitionFn, _clicked,
        _this = this;
      e.preventDefault();
      _clicked = this;
      if (dom.hasClass(this.overlay, "open")) {
        dom.removeClass(this.overlay, "open");
        dom.removeClass(this.container, "overlay-open");
        dom.addClass(this.overlay, "close");
        onEndTransitionFn = function(ev) {
          if (_this.support.transitions) {
            if (ev.propertyName !== "visibility") {
              return;
            }
            _this.overlay.removeEventListener(_this.transEndEventName, onEndTransitionFn);
          }
          dom.removeClass(_this.overlay, "close");
        };
        if (this.support.transitions) {
          this.overlay.addEventListener(this.transEndEventName, onEndTransitionFn);
        } else {
          onEndTransitionFn();
        }
      } else if (!dom.hasClass(this.overlay, "close")) {
        dom.addClass(this.overlay, "open");
        dom.addClass(this.container, "overlay-open");
      }
      /*
      if @overlay.hasClass("open")
        @overlay.removeClass "open"
        @container.removeClass "overlay-open"
        @overlay.addClass "close"
        onEndTransitionFn = (ev) =>
          if @support.transitions
            if ev.propertyName isnt "visibility" 
              return
            _clicked.off @transEndEventName, onEndTransitionFn
          @overlay.removeClass "close"
          return
        if @support.transitions
          @overlay.on @transEndEventName, onEndTransitionFn
        else
          onEndTransitionFn()
      else unless @overlay.hasClass("close")
        @overlay.addClass "open"
        @container.addClass "overlay-open"
      return
      */

    };

    return OverlayManager;

  })();

  /*
  domready( () ->
    getACar.start()
  )
  */


  $(function() {
    return getACar.start();
  });

}).call(this);
