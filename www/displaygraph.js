var G = new jsnx.Graph();
let loadFile = function(filePath) {
    var result = null;
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open("GET", filePath, false);
    xmlhttp.send();
    if (xmlhttp.status==200) {
      result = xmlhttp.responseText;
    }
    return result;
}

// var getDataFromRest = function() {
//     var self = this;
//         self.tasksURI = 'http://localhost:5000/giveme';
//         //self.tasks = ko.observableArray();

//         self.ajax = function(uri, method) {
//             var request = {
//                 url: uri,
//                 type: method,
//                 contentType: "application/json",
//                 "Access-Control-Allow-Origin": "*",
//                 accepts: "application/json",
//                 cache: false,
//                 dataType: 'json',
//             };
//             return $.ajax(request);
//         }

//         self.ajax(self.tasksURI, 'GET').done(function(data) {
//             console.log(data)
//         })

// }

var DisplayData = function () {
    var request = new XMLHttpRequest()

    request.open('GET', 'http://localhost:5000/giveme', true)
    request.onload = function() {
    // Begin accessing JSON data here
    var data = JSON.parse(this.response)

    if (request.status >= 200 && request.status < 400) {
            // var G = new jsnx.Graph();
            Object.keys(data).forEach(function(key) {
                //console.table('Key : ' + key + ', Value : ' + data[key])
                data[key].forEach( names => {
                    //console.log(names)
                    G.addEdge(key, names);
                })
            })
        drawMyGraph(G)
        }
    }
    request.send()
}

var addToGraph = function (nodeName, linkName) {
    G.addEdge(nodeName, linkName);
    drawMyGraph(G)
}
var sendData = function (nodeName, linkName) {
    var request = new XMLHttpRequest()
    request.open('POST', 'http://localhost:5000/addme', true)
    request.setRequestHeader("Content-Type", "application/json");
    request.setRequestHeader("Access-Control-Allow-Origin", "*");
    request.setRequestHeader("Access-Control-Allow-Methods", "GET, PUT, POST");
    request.onreadystatechange = function() {
        // Begin accessing JSON data here
        console.log(request.responseText)
        if (this.readyState == 4 && this.status == 200) {
            var data = JSON.parse(this.responseText)
            console.log("This is data : " + data.status)
            document.getElementById('replyMessage').innerHTML = data.status
        }
    }
    var network = JSON.stringify({ "nodeName": nodeName, "linkTo": linkName})
    request.send(network)
}

var drawMyGraph = function(G) {
    var color = d3.scale.category20();
    jsnx.draw(G, {
        element: '#canvas',
        weighted: true,
        height: 720,
        width: 1280,
        layoutAttr: {
            charge: -2000,
            linkDistance: 100
        },
        nodeAttr: {
            r: function(d) {
                // `d` has the properties `node`, `data` and `G`
                //console.log()
                return d.node.length * 5;
            },
            title: function(d) { return d.label;}
        },
        nodeStyle: {
            fill: function(d) { 
                return color(d.data.group); 
            },
            stroke: 'none'
        },
        edgeStyle: {
            fill: '#999',
            'stroke-width': 10
        },
        withLabels: true,
        labelStyle: {fill: 'white'},
    });
}