new Vue({
    el: '#app',

    props: ['username'],
    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        username: null, // Our username
        joined: false // True if email and username have been filled in
    },
    created: function() {

        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msgFull = JSON.parse(e.data);
            var msgType = msgFull.type;
            var msgBody = JSON.parse(atob(msgFull.Msg));
            switch (msgType) {
                case 'chat':
                    var element = document.getElementById('chat-messages');
                    self.chatContent += '<div class = "chip" id = "my-chip">'
                        + '<img src="' + self.gravatarURL(msgBody.email) + '">'
                        + msgBody.username
                        + '</div>';
                    if (msgBody.username === self.username) {
                        self.chatContent += '<font color="red">' + emojione.toImage(msgBody.message) + '</font><br/>';
                    } else {
                        self.chatContent += emojione.toImage(msgBody.message) + '<br/>';
                    }
                    element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
                    break;
                case 'count':
                    var element = document.getElementById('user-count');
                    var count = msgBody.count;
                    element.innerHTML = "Online #" + count;
                    break;
                default:
                    console.log("didn't recognize type: " + msgType);
            }
        });

    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() // Strip out html
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },
        join: function () {
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },
        gravatarURL: function(email) {
                return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
            }
        }
    }
);
