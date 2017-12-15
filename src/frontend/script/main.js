Math.logRand = function(min, max){
    var range = Math.abs(max - min);
    var scale = Math.ceil(Math.random() * Math.log10(range + 0.5));
    var value = Math.floor(Math.random() * (range + 1));
    return Math.round(value / Math.pow(10.0, scale)) + min;
}

Math.isOdd = function(num){
    return num % 2 ? true : false;
}

Math.collatzSequence = function(num){
    var sequence = []; sequence.push(num);
    while(num > 1){
        if(Math.isOdd(num)){ num = 3 * num + 1; }
        else{ num = num / 2; }
        sequence.push(num);
    } // console.log(sequence);
    return sequence;
}

$(document).ready(function(){

    // server url (ws/wss)
    var ws_url = 'ws://51.15.40.128:8080/ws' ; // test server: 'wss://echo.websocket.org'
    var ws; // global scope for websocket instance
    
    var $input = $('#collatzNumber'); // field with a number
    var $button = $('#collatzAction'); // two functional button
    var $output = $('#collatzOutput'); // place to output results
    var $results = $('#collatzResults'); // results block
    var $random = $('#collatzRandom'); // random button
    var $example = $('#collatzExample'); // example button
    var $debug = $('#debugOutput'); // for messages
    
    var maximum = 4 * 8 * 15 * 16 * 23 * 42 * 108; // DHARMA INITIATIVE
    var minimum = 3; // 1 is done, 2 is only 1 step, 3 is okay.
    var counter = 0; // answer counter
    
    function addAnswer(number, data){
        var result, elements = [];
        try {
            json = JSON.parse(data);
            elements.push('Number: ' + json['Number']);
            elements.push('Length: ' + json['Path'].length); // fix: get from server
            elements.push('Highest: ' + json['Path'].reduce(function(a,b){return a>b?a:b})); // fix: get from server
            elements.push('Average: ' + json['Path'].reduce(function(a,b){return a+b}) / json['Path'].length); // fix: get from server
            elements.push('Time: ' + (json['Time']/1000) + 'ms');
            result = elements.join(', ');
        } catch (e) { console.log(e); result = data; }
        result = 'Answer #' + number + ': ' + result;
        var item = '<div class="item">' + result + '</div>';
        $output.append(item);
        $output.scrollTop($output[0].scrollHeight);
    }
    
    function logDebug(data){
        var item = data + '\n';
        $debug.append(item);
        $debug.animate({scrollTop: $debug[0].scrollHeight}, 200, 'swing');
    }
    
    $random.on('click', function(ev){
        $input.val(Math.logRand(minimum, maximum));
    });
    
    $example.on('click', function(ev){
        var value = $input.val();
        var result = Math.collatzSequence(value);
        logDebug('Path for '+ value +' is ' + result.join(' > '));
    });
    
    $input.val(Math.logRand(minimum, maximum));
    
    $button.attr('data-opened', 'false');
    
    $button.on('click', function(ev){
        
        if($button.attr('data-opened') == 'false'){
            
            $button.attr('disabled','disabled');
            $button.val('Connecting...');
                
            $output.html('');
            $debug.html('');
            
            ws = new WebSocket(ws_url); // new connection
            
            counter = 0;
            
            ws.onopen = function(){
                console.log('ws opened');
                logDebug('Connection opened');
                
                $button.removeAttr('disabled');
                $button.attr('data-opened', 'true');
                $button.val('Stop');
                
                var value = $input.val();
                ws.send(value);
                logDebug('Message sent: ' + value);
            };
            
            ws.onmessage = function(ev){
                console.log('ws answer'); // console.log('ws answer: ' + ev.data);
                addAnswer(counter++, ev.data);
                if(counter == 1) $results.slideDown(300);
            };
            
            ws.onerror = function(){
                console.log('ws error');
                logDebug('Error');
            };
            
            ws.onclose = function(ev){
                console.log('ws closed: ' + ev.code);
                logDebug('Connection closed (' + ev.code + ')');
                
                $input.val(Math.logRand(minimum, maximum));
                
                $button.attr('data-opened', 'false');
                $button.removeAttr('disabled');
                $button.val('Start');
            };
            
        } else {
            
            $button.attr('disabled','disabled');
            $button.val('Disconnecting...');
            
            ws.close();
            
        }
        
    });
    
});
