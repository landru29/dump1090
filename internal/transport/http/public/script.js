const xmlhttp = new XMLHttpRequest();
const aircraftDict = {};
const aircraftMarkerIcon = L.icon({
	iconUrl: "ac.png",
	iconSize: [64,64],
	iconAnchor: [32,32]
});

const selectedAircraftMarkerIcon = L.icon({
	iconUrl: "acs.png",
	iconSize: [64,64],
	iconAnchor: [32,32]
});

window.addEventListener('load', function() {
    const map = L.map('map').setView([48.12, -1.86], 9);

	const tiles = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
		maxZoom: 19,
		minZoom: 3,
		attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
	}).addTo(map);

	const acContainer = document.getElementById("ac");

	aircraftGetter = function() {
		loadAircrafts().then(function(aircraftList) {
			aircraftList.forEach((elt) => {
				elt.seen = new Date();
				const oldAC = aircraftDict[elt.icao];
	
				
				if (oldAC != undefined) {
					aircraftDict[elt.icao] = Object.assign(oldAC, elt);
				} else {
					aircraftDict[elt.icao] = elt;
					const container = document.createElement('div');
					container.className = 'aircraft';
					container.id = `icao_${elt.icao}`;
					aircraftDict[elt.icao].container = container;
					acContainer.appendChild(container);
					container.addEventListener('click', () => clickAircraft(map, elt.icao));
				}
	
				if (aircraftDict[elt.icao].lat != 0 && aircraftDict[elt.icao].lon != 0) {
					const coord = new L.LatLng(aircraftDict[elt.icao].lat, aircraftDict[elt.icao].lon);
					aircraftDict[elt.icao].coordinate = coord;
					if (aircraftDict[elt.icao].marker == undefined) {
						aircraftDict[elt.icao].marker = L.marker(coord, {icon: aircraftMarkerIcon, icao: elt.icao}).addTo(map);
						aircraftDict[elt.icao].marker.on('click', ()=>{
							clickAircraft(map, elt.icao);
						});
					}
					
					aircraftDict[elt.icao].marker.setLatLng(coord); 
					aircraftDict[elt.icao].marker.setRotationAngle(elt.track);
				}
			});
			
			// clean outdated aircrafts.
			Object.keys(aircraftDict).forEach((key)=>{
				if (aircraftDict[key].seen.getTime() < new Date().getTime()- 120000) {
					acContainer.removeChild(aircraftDict[key].container);
					delete aircraftDict[key];
				}
			});
	
			// refresh aircraft data in the table.
			Object.keys(aircraftDict).forEach((key)=>{
				const icao = document.createElement('span');
				icao.className='icao';
				icao.innerHTML=aircraftDict[key].icao_hex;
	
				const flight = document.createElement('span');
				flight.className='flight';
				flight.innerHTML=aircraftDict[key].flight;
	
				const altitude = document.createElement('span');
				altitude.className='altitude';
				altitude.innerHTML=`${aircraftDict[key].altitude}ft`;
	
				const speed = document.createElement('span');
				speed.className='speed';
				speed.innerHTML=`${aircraftDict[key].speed}kt`;
	
				aircraftDict[key].container.innerHTML='';
				aircraftDict[key].container.appendChild(icao);
				aircraftDict[key].container.appendChild(flight);
				aircraftDict[key].container.appendChild(altitude);
				aircraftDict[key].container.appendChild(speed);
			});
		});
	}

	aircraftGetter();

	this.setInterval(aircraftGetter(), 2000);
});

function clickAircraft(map, icao) {
	const line = document.getElementById(`icao_${icao}`);
	const isHighlight = (hasClass(line, 'highlight'));
	
	for (const child of document.getElementById("ac").children) {
		child.className = 'aircraft'
	}

	Object.keys(aircraftDict).forEach((ac) => aircraftDict[ac].marker.setIcon(aircraftMarkerIcon));

	if (!isHighlight) {
		line.className = 'aircraft highlight';
		aircraftDict[icao].marker.setIcon(selectedAircraftMarkerIcon);

		map.panTo(aircraftDict[icao].coordinate);
	}
}

function hasClass(elt, className) {
	return(elt.className.split(' ').filter((name) => {
		return name === className;
	}).length !== 0);
}


function loadAircrafts() {
	return new Promise((resolve, reject) => {
		xmlhttp.open('GET', window.apiPath, true);
		xmlhttp.onreadystatechange = () => {
			if (xmlhttp.readyState == 4) {
				if(Math.floor(xmlhttp.status/100) == 2) {
					resolve(JSON.parse(xmlhttp.responseText));
				} else {
					reject({
						code: xmlhttp.status,
						msg: xmlhttp.responseText,
					});
				}
			}
		};
		xmlhttp.send(null);
	});
}