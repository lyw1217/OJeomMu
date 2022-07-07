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
    if (rad[0]['value'] == 0 || cat[0]['value'] == "none") {
        alert("반경과 카테고리를 선택해주세요.", "", "info")
        return
    }
    //var homereturn = $('#homereturn').is(":checked")

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
        contentType: "application/json; charset=utf-8",
        url: '/sendToGo',
        data: JSON.stringify(params),
        error: function () {
            alert("조회 중 에러가 발생했어요..");
        },
        success: function (resData) {
            if (!$.trim(resData)) {
                alert("주변에 음식점이 없어요.", "", "error");
            } else {
                let msg = "";               
                if (resData.CAT_NAME) {
                    msg = msg + "카테고리\t: " + resData.CAT_NAME + "\n";
                } 
                if (resData.DISTANCE != 0) {
                    msg = msg + "거리\t: " + resData.DISTANCE + " 미터\n";
                }
                if (resData.ROAD_ADDRESS) {
                    msg = msg + "주소\t: " + resData.ROAD_ADDRESS + "\n";
                }
                /*
                if (resData.URL) {
                    msg = msg + "URL\t: " + resData.URL + "\n";
                }
                */
                if (resData.PHONE) {
                    msg = msg + "전화번호\t: " + resData.PHONE;
                }
                
                // 검색된 음식점의 좌표로 마커 및 지도 이동
                setMarkerPosition(resData.Y, resData.X);

                //resultAlert(resData.NAME, msg, resData.URL, "success", homereturn, params['x'], params['y'] );
                resultAlert(resData.NAME, msg, resData.URL, "success" );
            }
        }
    });
}

var alert = function(title, msg, icon) {
    swal({
        title : title,
        text : msg,
        icon : icon,
        time : 1500,
        showConfirmButton : false
    });
}

//var resultAlert = function(title, msg, url, icon, hm, x, y) {
var resultAlert = function(title, msg, url, icon) {
    
    swal({
        title : title,
        text : msg,
        icon : icon,
        buttons: {
            cancel: "시러", 
            catch: {
                text: "가자!",
                value: "go",
        }},
        showConfirmButton : true
    })
    .then((value) => {
        switch (value) {
            case "go":
                if (isValidHttpUrl(url)) {
                    window.location.href=url;
                } else {
                    swal("URL이 올바르지 않아요.", {
                        icon : "warning"
                    })
                }
                break;
            default:
                /*
                if ( hm ) {
                    // 원위치 복귀
                    setMarkerPosition(y, x);
                }
                */
        }
    });
}

// https://stackoverflow.com/questions/5717093/check-if-a-javascript-string-is-a-url
function isValidHttpUrl(string) {
    let url;

    try {
        url = new URL(string);
    } catch (_) {
        return false;  
    }

    return url.protocol === "http:" || url.protocol === "https:";
}