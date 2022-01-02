/*!
    * Start Bootstrap - SB Admin v7.0.4 (https://startbootstrap.com/template/sb-admin)
    * Copyright 2013-2021 Start Bootstrap
    * Licensed under MIT (https://github.com/StartBootstrap/startbootstrap-sb-admin/blob/master/LICENSE)
    */
// 
// Scripts
// 

window.addEventListener('DOMContentLoaded', event => {

    // Toggle the side navigation
    const sidebarToggle = document.body.querySelector('#sidebarToggle');
    if (sidebarToggle) {
        // Uncomment Below to persist sidebar toggle between refreshes
        // if (localStorage.getItem('sb|sidebar-toggle') === 'true') {
        //     document.body.classList.toggle('sb-sidenav-toggled');
        // }
        sidebarToggle.addEventListener('click', event => {
            event.preventDefault();
            document.body.classList.toggle('sb-sidenav-toggled');
            localStorage.setItem('sb|sidebar-toggle', document.body.classList.contains('sb-sidenav-toggled'));
        });
    }

});


/* GIN <-> Ajax */
function sendToGo() {
    var rad = $('#radius').serializeArray();
    var cat = $('#category').serializeArray();

    loc = marker.getPosition();

    //json 가공
    //[{ name : "a", value : "1" }] to {"a":"1"}
    var params = {};
    params[rad[0]['name']] = rad[0]['value'];
    params[cat[0]['name']] = cat[0]['value'];
    params['x'] = String(loc.getLng());
    params['y'] = String(loc.getLat());
    
    $.ajax({
        type: 'POST',
        url: '/sendToGo',
        data: JSON.stringify(params),
        //dataType : 'json',
        //contentType : "application/json; charset=UTF-8",
        error: function () {
            alert("에러 발생");
        },
        success: function (json) {
            //alert(json)
        }
    });

}
