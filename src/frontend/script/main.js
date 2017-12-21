Math.logRand = function(min, max){
    var range = Math.abs(max - min);
    var scale = Math.ceil(Math.random() * Math.log10(range + 0.5));
    var value = Math.floor(Math.random() * (range + 1));
    return Math.round(value / Math.pow(10.0, scale)) + min;
}

Math.isOdd = function(num){
    return num % 2 ? true : false;
} // Math.isEven = function(num){ return num % 2 ? false : true;} 

Math.collatzSequence = function(num){
    var sequence = []; sequence.push(num);
    while(num > 1){
        if(Math.isOdd(num)){ num = 3 * num + 1; }
        else{ num = num / 2; }
        sequence.push(num);
    } // console.log(sequence);
    return sequence;
}

// when page is ready
$(document).ready(function(){

    // server url (ws/wss)
    var ws_url = 'wss://echo.websocket.org' ; // test server: 'wss://echo.websocket.org'
    var ws; // global scope for websocket instance
    
    var test_json = '{"Number":1,"PathLength":1,"MaxNumber":1,"AverageNumber":1,"Time":1}';
    
    var $input = $('#collatzNumber'); // field with a number
    var $button = $('#collatzAction'); // two functional button
    
    var $debug = $('#debugOutput'); // for messages
    var $random = $('#collatzRandom'); // random button
    
    var $exampleButton = $('#collatzExampleButton'); // example button
    var $exampleResult = $('#collatzExampleResult'); // example button
    
    var $results = $('#collatzResults'); // results block
    var $answers = $('#collatzOutput'); // place to output results
    var $histogram = $('#collatzHistogram'); // place for histogram
    
    var counter = 0; // answer counter
    
    var maximum = 4 * 8 * 15 * 16 * 23 * 42 * 108; // DHARMA INITIATIVE
    var minimum = 3; // 1 is done, 2 is only 1 step, 3 is okay.

    // for collecting stats
    Histogram = {
        values: [0], total: 0, maximum: 0,
        reset: function(){
            this.values = [0];
            this.total = 0;
            this.maximum = 0;
        },
        add: function(value){
            if(this.values.length < value + 1){
                for(var i = this.values.length; i < value + 1; i++){
                    this.values.push(0);
                }
            }
            this.values[value]++; this.total++;
            if(this.values[value]>this.maximum){
                this.maximum = this.values[value];
            }
        },
        print: function(){
            var result = [];
            for(var i = 0; i < this.values.length; i++){
                if(this.values[i]>0){
                    result.push(i+':'+this.values[i]);
                }
            }
            return result.join(', ');
        }
    };
    
    function createDiv(content, id) {
        var elem = document.createElement('div');
        elem.classList.add("item"); elem.id = id;
        elem.appendChild(document.createTextNode(content));
        return elem;
    }
    
    function addAnswer(data){
        counter++;
        var result, maximum = 100, update = 10;
        
        if (typeof data === 'object'){
            
            Histogram.add(data['PathLength']);
            
            var elements = []; // prepare string (#GOVNOKOD)
            elements.push('Number: ' + data['Number']);
            elements.push('Length: ' + data['PathLength']);
            elements.push('Highest: ' + data['MaxNumber']);
            elements.push('Average: ' + data['AverageNumber']);
            elements.push('Time: ' + (data['Time']/1000) + 'ms'); // is it really ms?
            result = elements.join(', ');
            
        } else {
            result = data;
            
            Histogram.add(Math.logRand(1,666)); // testing
            // console.log(Histogram.values);
        }
        
        if(Histogram.total % update == 0){
            // $histogram.text(Histogram.print());
            Plotly.relayout($histogram[0], {y: Histogram.values});
        }
            
        $answers.queue(function(){ // experimental sh!t (that works)
            // limit shown elements by removing oldest
            if(counter>=maximum) {
                $(this).children('#answer_'+(counter-maximum)).queue(function(){
                    $(this).remove().dequeue();
                });
            }
            
            // add result
            // var item = '<div id="answer_'+counter+'" class="item">' + result + '</div>'; // create div instead?
            var item = createDiv(result, 'answer_'+counter);
            $(this).append(item).ready(function(){
                item.scrollIntoView({block: "end", behavior: "smooth"}); // smooth by css
            });
                
            $(this).dequeue();
        });
    }
    
    function logDebug(data){
        $debug.append(data+'\n').animate({scrollTop: $debug[0].scrollHeight}, 200, 'swing');
    }
    
    $button.attr('data-opened', 'false');
    
    $random.on('click', function(ev){ $input.val(Math.logRand(minimum, maximum)); });
    
    $exampleButton.on('click', function(ev){
        var value = $input.val();
        var result = Math.collatzSequence(value);
        $exampleResult.text('Path for '+ value +' is ' + result.join(' > '));
    });
    
    $button.on('click', function(ev){
        
        if($button.attr('data-opened') == 'false'){
            
            $button.attr('disabled','disabled').val('Connecting...');
            
            $debug.html('');
            
            $answers.html('');
            
            Histogram.reset();
            
            // histo init
            Plotly.plot($histogram[0], [{y: Histogram.values, type: 'bar'}], {showlegend: false}, {staticPlot: true});
            
            $results.slideDown(500);
            
            //* for testing
            var test = setInterval(function(){
                var rand = Math.random()*3;
                for(var i = 0; i < rand; i++){
                    addAnswer('test '+counter);
                }
            },200);
            //*/
            
            logDebug('Connection attempt');
            ws = new WebSocket(ws_url); // new connection
            
            counter = 0;
            
            ws.onopen = function(){
                // ws.send(test_json);
                console.log('ws opened');
                logDebug('Connection is opened');
                
                $button.attr('data-opened', 'true').removeAttr('disabled').val('Stop');
                
                var value = $input.val();
                
                ws.send(value);
                logDebug('Message sent: "' + value + '"');
            };
            
            ws.onmessage = function(ev){
                console.log('ws answer'); // if(counter==0) console.log(ev);
                addAnswer($.parseJSON(ev.data));
            };
            
            ws.onerror = function(){
                console.log('ws error');
                logDebug('An error has occurred');
            };
            
            ws.onclose = function(ev){
                console.log('ws closed: ' + ev.code);
                logDebug('Connection is closed (' + ev.code + ')');
                
                $button.attr('data-opened', 'false').removeAttr('disabled').val('Start');
                // $input.val(Math.logRand(minimum, maximum)); // UX =)
                
                clearInterval(test); // testing
            };
            
        } else {
            
            $button.attr('disabled','disabled').val('Disconnecting...');
            
            ws.close();
            
        }
        
    });
    
    $input.val(Math.logRand(minimum, maximum));
    
});
