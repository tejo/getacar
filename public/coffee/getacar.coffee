dom = {
  hasClass : (elem, className) ->
    return new RegExp(' ' + className + ' ').test(' ' + elem.className + ' ')
  addClass : (elem, className) ->
    if (!dom.hasClass(elem, className))
      elem.className += ' ' + className
  removeClass : (elem, className) ->
    newClass = ' ' + elem.className.replace( /[\t\r\n]/g, ' ') + ' '
    if (dom.hasClass(elem, className))
      while (newClass.indexOf(' ' + className + ' ') >= 0 )
        newClass = newClass.replace(' ' + className + ' ', ' ')
      elem.className = newClass.replace(/^\s+|\s+$/g, '')
};

getACar = {
  start : ->
    overLayManager = new OverlayManager();
    FastClick.attach(document.body);
};

class OverlayManager
  constructor: () ->
    @container = document.querySelector( 'div.container' )#$("div.container")
    @triggerBttn = document.getElementById( 'trigger-overlay' )#$("#trigger-overlay")
    @overlay = document.querySelector( 'div.overlay' )#$("div.overlay")
    @closeBttn = @overlay.querySelector( 'button.overlay-close' )#@overlay.find("button.overlay-close")
    @transEndEventNames =
      WebkitTransition: "webkitTransitionEnd"
      MozTransition: "transitionend"
      OTransition: "oTransitionEnd"
      msTransition: "MSTransitionEnd"
      transition: "transitionend"

    @transEndEventName = @transEndEventNames[Modernizr.prefixed("transition")]
    @support = transitions: Modernizr.csstransitions

    #@triggerBttn.on "click", @toggleOverlay
    #@closeBttn.on "click", @toggleOverlay
    @triggerBttn.addEventListener( 'click', @toggleOverlay )
    @closeBttn.addEventListener( 'click', @toggleOverlay )

  toggleOverlay : (e) =>
    e.preventDefault();
    _clicked  = this
    if dom.hasClass(@overlay,"open")
      dom.removeClass(@overlay, "open")
      dom.removeClass(@container, "overlay-open")
      dom.addClass(@overlay, "close")
      onEndTransitionFn = (ev) =>
        if @support.transitions
          if ev.propertyName isnt "visibility" 
            return
          @overlay.removeEventListener @transEndEventName, onEndTransitionFn
        dom.removeClass(@overlay, "close")
        return
      if @support.transitions
        @overlay.addEventListener @transEndEventName, onEndTransitionFn
      else
        onEndTransitionFn()
    else unless dom.hasClass(@overlay,"close")
      dom.addClass(@overlay, "open")
      dom.addClass(@container, "overlay-open")
    return
    ###
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
    ###




# Dom Ready
###
domready( () ->
  getACar.start()
)
###
$ ->
  getACar.start()