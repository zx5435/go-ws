'use strict'
var socket;

$(function() {
    socket = io(':5000');

    socket.emit('testcallback', 'copy2', function(data) {
        console.log(data)
    });

    socket.on('new msg', function(data) {
        document.writeln(data.username, ' ', data.message, '<br>')
    });
});

function send(txt) {
    socket.emit('new msg', txt);
}