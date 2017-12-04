Math.logRand = function(min, max){
    
    var range = Math.abs(max - min);
    
    var scale = Math.ceil(Math.random() * Math.log10(range + 0.5));
    
    var value = Math.floor(Math.random() * (range + 1));
    
    return Math.round(value / Math.pow(10.0, scale)) + min;
    
}

function isOdd(num){
    return num % 2;
}

function collatzSequence(num){
    var sequence = [];
    while(num > 1){
        if(isOdd(num)){ num = 3 * num + 1; }
        else{ num = num / 2; }
        sequence.push(num);
    } // console.log(sequence);
    return sequence;
}

$(document).ready(function(){

    // server url (ws/wss)
    var ws_url = 'wss://echo.websocket.org'; // test server
    var ws; // global scope for websocket instance
    
    var $input = $('#collatzNumber'); // field with a number
    var $button = $('#collatzAction'); // two functional button
    var $output = $('#collatzOutput'); // place to output results
    var $results = $('#collatzResults'); // results block
    
    var maximum = 4 * 8 * 15 * 16 * 23 * 42 * 108;
    var minimum = 3; // 1 is done, 2 is only 1 step, 3 is okay. 
    
    $input.val(Math.logRand(minimum, maximum));
    
    $button.attr('data-opened', 'false');
    
    $output.text('Nothing yet...');
    
    function addResult(data){
        
        var item = '<div class="item">' + data + '</div>';
        
        $output.append(item);
        
    }
    
    $button.on('click', function(ev){
        
        var counter = 0;
        
        if($button.attr('data-opened') == 'false'){
            
            $results.slideDown(300);
            
            $button.attr('disabled','disabled');
            $button.val('Connecting...');
                
            $output.html('');
            
            ws = new WebSocket(ws_url); // new connection
            
            ws.onopen = function(){
                
                console.log('ws opened');
                
                $button.removeAttr('disabled');
                $button.attr('data-opened', 'true');
                $button.val('Stop');
                
                addResult('Connection opened');
                
                var value = $input.val();
                
                ws.send(value);
                addResult('Integet sent: ' + value);
                
                addResult('Awaited result: ' + value + ', ' + collatzSequence(value).join(", "));
                
            };
            
            ws.onmessage = function(ev){
                console.log('ws answer: ' + ev.data);
                counter++;
                addResult('Answer #' + counter + ': ' + ev.data);
            };
            
            ws.onerror = function(){
                console.log('ws error');
                $button.val('Error');
            };
            
            ws.onclose = function(ev){
                
                console.log('ws closed: ' + ev.code);
                
                addResult('Connection closed (' + ev.code + ')');
                
                $button.attr('data-opened', 'false');
                $button.removeAttr('disabled');
                $button.val('Start');
                
            };
            
        } else {
            
            $button.attr('disabled','disabled');
            $button.val('Disconnecting...');
            
            ws.close();
            
            $input.val(Math.logRand(minimum, maximum));
            
        }
        
    });
    
});
