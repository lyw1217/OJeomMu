// kakaomap scripts

var mapContainer = document.getElementById("map"), // 지도를 표시할 div
  mapOption = {
    center: new kakao.maps.LatLng(33.450701, 126.570667), // 지도의 중심좌표
    level: 2, // 지도의 확대 레벨
    mapTypeId: kakao.maps.MapTypeId.ROADMAP, // 지도종류
  };

// 지도를 생성한다
var map = new kakao.maps.Map(mapContainer, mapOption);

// 주소-좌표 변환 객체를 생성합니다
var geocoder = new kakao.maps.services.Geocoder();

// 마커 생성
var marker = new kakao.maps.Marker();
// 디폴트 마커 생성
var latlng = new kakao.maps.LatLng(33.450701, 126.570667);
marker.setPosition(latlng);
marker.setMap(map);
// 일반 지도와 스카이뷰로 지도 타입을 전환할 수 있는 지도타입 컨트롤을 생성합니다
var mapTypeControl = new kakao.maps.MapTypeControl();

// HTML5의 geolocation으로 사용할 수 있는지 확인합니다
if (navigator.geolocation) {
  // GeoLocation을 이용해서 접속 위치를 얻어옵니다
  navigator.geolocation.getCurrentPosition(function (position) {
    var lat = position.coords.latitude, // 위도
      lon = position.coords.longitude; // 경도

    var locPosition = new kakao.maps.LatLng(lat, lon); // 마커가 표시될 위치를 geolocation으로 얻어온 좌표로 생성합니다

    searchDetailAddrFromCoords(locPosition, function (result, status) {
        if (status === kakao.maps.services.Status.OK) {
            if ( !!result[0].road_address ) {
                var detailAddr = "<div>" 
                + result[0].road_address.address_name
                + " "
                + result[0].road_address.building_name 
                + "</div>";
            } else {
                var detailAddr = !!result[0].address
                ? "<div>" 
                + result[0].address.address_name
                + "</div>"
                : "<div><br></div>";
            }

          // 마커를 클릭한 위치에 표시합니다
          marker.setPosition(locPosition);
          marker.setMap(map);
    
          displayInfo(detailAddr, status);

          // 지도 중심좌표를 접속위치로 변경합니다
          map.setCenter(locPosition);
        }
      });
  });
} else {
  // HTML5의 GeoLocation을 사용할 수 없을때 마커 표시 위치와 인포윈도우 내용을 설정합니다

  var locPosition = new kakao.maps.LatLng(33.450701, 126.570667),
    message = "geolocation을 사용할수 없어요..";

  displayMarker(locPosition);
}

// 지도에 마커와 인포윈도우를 표시하는 함수입니다
function displayMarker(locPosition) {
  marker.setPosition(locPosition);

  marker.setDraggable(false);
  
  // 기존에 마커가 있다면 제거
  marker.setMap(null);
  marker.setMap(map);

  // 지도 중심좌표를 접속위치로 변경합니다
  map.setCenter(locPosition);
}

// 지도에 컨트롤을 추가해야 지도위에 표시됩니다
// kakao.maps.ControlPosition은 컨트롤이 표시될 위치를 정의하는데 TOPRIGHT는 오른쪽 위를 의미합니다
map.addControl(mapTypeControl, kakao.maps.ControlPosition.TOPRIGHT);

// 지도 확대 축소를 제어할 수 있는  줌 컨트롤을 생성합니다
var zoomControl = new kakao.maps.ZoomControl();
map.addControl(zoomControl, kakao.maps.ControlPosition.RIGHT);

// 지도에 클릭 이벤트를 등록합니다
// 지도를 클릭하면 마지막 파라미터로 넘어온 함수를 호출합니다
kakao.maps.event.addListener(map, "click", function (mouseEvent) {
  searchDetailAddrFromCoords(mouseEvent.latLng, function (result, status) {
    if (status === kakao.maps.services.Status.OK) {
        if ( !!result[0].road_address ) {
            var detailAddr = "<div>" 
            + result[0].road_address.address_name
            + " "
            + result[0].road_address.building_name 
            + "</div>";
        } else {
            var detailAddr = !!result[0].address
            ? "<div>" 
            + result[0].address.address_name
            + "</div>"
            : "<div><br></div>";
        }
        
      // 마커를 클릭한 위치에 표시합니다
      marker.setPosition(mouseEvent.latLng);
      marker.setMap(map);

      displayInfo(detailAddr, status);
    }
  });
});

// 마커에 대한 주소정보를 표출하는 함수입니다
function displayInfo(result, status) {
  if (status === kakao.maps.services.Status.OK) {
    var infoDiv = document.getElementById("currAddr");

    infoDiv.innerHTML = result;
  }
}

function searchAddrFromCoords(coords, callback) {
  // 좌표로 행정동 주소 정보를 요청합니다
  geocoder.coord2RegionCode(coords.getLng(), coords.getLat(), callback);
}

function searchDetailAddrFromCoords(coords, callback) {
  // 좌표로 법정동 상세 주소 정보를 요청합니다
  geocoder.coord2Address(coords.getLng(), coords.getLat(), callback);
}

function locationLoadSuccess(pos) {
  // 현재 위치 받아오기
  var currentPos = new kakao.maps.LatLng(
    pos.coords.latitude,
    pos.coords.longitude
  );

  // 지도 이동(기존 위치와 가깝다면 부드럽게 이동)
  map.panTo(currentPos);

  // 마커 생성
  marker.setPosition(currentPos);

  // 기존에 마커가 있다면 제거
  marker.setMap(null);
  marker.setMap(map);
}

function locationLoadError(pos) {
  alert("위치 정보를 가져오는데 실패했습니다.");
}

// 위치 가져오기 버튼 클릭시
function getCurrentPosBtn() {
  navigator.geolocation.getCurrentPosition(
    locationLoadSuccess,
    locationLoadError
  );
}

function setStorePosition(lat, lng) {

  var latlng = new kakao.maps.LatLng(lat, lng);
  
  map.panTo(latlng);

  marker.setPosition(latlng);
  // 기존에 마커가 있다면 제거
  marker.setMap(null);
  marker.setMap(map);
}