Math.logRand = function(min, max){
    var range = Math.abs(max - min);
    var scale = Math.round(Math.random() * Math.log10(range + 0.5));
    var value = Math.floor(Math.random() * (range + 1));
    return Math.round(value / Math.pow(10.0, scale)) + min;
}

// Math.isEven = function(num){ return num % 2 ? false : true;} 
Math.isOdd = function(num){ return num % 2 ? true : false; }

Math.collatzSequence = function(num){
    var sequence = [num];
    while(num > 1){
        if(Math.isOdd(num)){ num = 3 * num + 1; }
        else{ num = num / 2; }
        sequence.push(num);
    } 
    return sequence;
}

// when page is ready
$(document).ready(function(){

    // default ws address if front launched from file
    var def_host = 'collatz.devheap.org';
    var def_proto = 'wss:';
    
    // host and protocol
    var ws_host, ws_proto;
    switch(location.protocol){
        case 'https:': ws_host = location.host; ws_proto = 'wss:'; break;
        case 'http:': ws_host = location.host; ws_proto = 'ws:'; break;
        case 'file:': default: ws_host = def_host; ws_proto = def_proto; break;
    }
    
    // ws server url
    var ws_url = ws_proto + '//' + ws_host + '/ws'; // test server: 'wss://echo.websocket.org'
    
    var ws; // global scope for websocket instance
    
    var n_minimum = 1; // collatz min
    var n_maximum = 999999999999999; // js max
    
    var max_answers = 300; // maximum answers shown at a time
    
    var $input = $('#collatzNumber'); // field with a number
    var $button = $('#collatzAction'); // two functional button
    
    var $debug = $('#collatzDebug'); // for messages
    var $random = $('#collatzRandom'); // random button
    
    var $results = $('#collatzResults'); // results block
    var $histogram = $('#collatzHistogram'); // place for histogram
    var $answers = $('#collatzOutput'); // place to output results
    
    var $example = $('#collatzExample'); // example placeholder

    // Histogram singleton for collecting stats
    Histogram = {
        values: [0], total: 0, maximum: 0, time: 0,
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
            return this.values.reduce(function(a, b, i){
                return a + ', ' + i + ':' + b;
            });
        },
        outdated: function(sec){
            if((Date.now() - this.time) > (sec * 1000)){
                this.time = Date.now();
                return true;
            } else return false;
        },
    };
    
    Answers = {
        maximum: 0, element: $answers,
        itemClass: 'item', itemId: 'item', itemUpdated: 'filled',
        add: function(id, content){
            if(id > this.maximum){ this.blocks(id); }
            $('#'+this.itemId+'_'+id).html(content).addClass(this.itemUpdated);
        },
        blocks: function(number){
            if(this.element != null){
                for(var i=this.maximum; i<number; i++){
                    var item = '<div class="'+this.itemClass+'" id="'+this.itemId+'_'+(i+1)+'">...</div>';
                    this.element.append(item);
                }
                this.maximum = number;
            }
        },
        format: function(obj){
            var fields = [obj['Number'], obj['PathLength'], obj['MaxNumber'], obj['AverageNumber'], (obj['Time']/1000000)+' ms'];
            return "<span>" + fields.join('</span><span>') + "</span>";
        }
    };
    
    // Scroll object (aka kostil)
    Scroll = {
        height: 0, position: 0, target: 0,
        interval: 1000 * 1/3, time: 0, easing: 3/5, // more magic numbers
        moving: function(){
            return (this.position < this.target);
        },
        ready: function(){
            return (Date.now() - this.time) > (this.interval);
        },
        update: function(el){
            this.height = el.scrollHeight;
            this.position = el.scrollTop + el.clientHeight;
            this.target = this.position + Math.floor((this.easing) * (this.height-this.position));
        },
        smooth: function(element){
            this.update(element[0]);
            if(this.moving() && this.ready()){
                element.animate({
                    scrollTop: (this.target - element.height() - 1)
                }, this.interval - 10, 'linear');
                this.time = Date.now();
            }
        },
    };
    
    // adds answer
    function addAnswer(data){
        
        var result = '', itemid = '';
        
        // prepare string
        if (typeof data === 'object'){
            result = Answers.format(data);
            itemid = data['Number'];
        } else { result = data; }
        
        // queue
        $answers.queue(function(){
        
            // limit shown elements by removing oldest
            var size = $answers.children().length;
            if((size > max_answers) && (size % Math.floor(max_answers / 3) == 0)){ // let 1/3 more of max to stay
                $answers.children().remove(":lt(" + (size - max_answers) + ")");
                Scroll.update($answers[0]); // better smooth scroll
            }
            
            // add result
            Answers.add(itemid, result);
            
            // scroll smoothly
            Scroll.smooth($answers);
            
            // next please
            $(this).dequeue();
            
        });
        
    }
    
    function updateHistogramView(interval){
        if(Histogram.outdated(interval)){
            var rand = Math.round(Math.random() * 100);
            setTimeout(function(){
                Plotly.restyle($histogram[0], 'y', [Histogram.values]);
                // still don't know which is better (restyle vs relayout).
                // Plotly.relayout($histogram[0], {y: Histogram.values});
            }, rand); // if it is not async then i am dumb
        }
    }
    
    // updates histogram
    function updateHistogram(data){
        if (typeof data === 'object'){ Histogram.add(data['PathLength']); }
        updateHistogramView(2);
    }
    
    // onpage log
    function logDebug(data){
        var dt = new Date();
        var message = dt.getHours() + ':' + dt.getMinutes()+ ':' + dt.getSeconds() + ' - ' + data + '\n';
        // $debug.append(message).animate({scrollTop: $debug[0].scrollHeight}, 200, 'swing');
        console.log(message);
    }
    
    function doAction(){
        
        // show results block
        $results.slideDown(500, function(){
            // show debug output
            // $debug.slideDown(500); 
        });
        
        if($button.attr('data-opened') == 'false'){
            
            $button.attr('disabled','disabled').val('Connecting...');
            
            // clear output
            $answers.children().remove();
            $answers.addClass('autoscroll');
            
            // histogram data init
            Histogram.reset();
            
            // histogram view style
            var margin = 40; // histogram border margin in px
            var font = { family: "Cuprum, sans-serif", size: 15, color: "#1f1f1f" };
            var layout = {
                margin: { l: margin, r: margin, b: margin, t: margin, pad: margin/8 },
                xaxis: { title : "Path length", titlefont: font, tickfont: font },
                yaxis: { title : "Count", titlefont: font, tickfont: font },
                showlegend: false,
            };
            // plot
            Plotly.newPlot($histogram[0], [{y: Histogram.values, type: 'bar'}], layout, {staticPlot: true});
            $(window).on('resize', function(){ Plotly.Plots.resize($histogram[0]); });
            
            // scroll to results
            $('html, body').animate({ scrollTop: $results.offset().top }, 500);
            
            // new connection
            try{
                ws = new WebSocket(ws_url);
                logDebug('Connection attempt');
            } catch (e) {
                logDebug('Fail to connect: ' + e.name + "; " + e.message + "\n" + e.stack); 
            }
            
            ws.onopen = function(){
                logDebug('Connection is opened');
                
                var val = $input.val();
                ws.send(val);
                logDebug('Message sent: ' + val);
                Answers.maximum = val;
            
                $button.attr('data-opened', 'true').removeAttr('disabled').val('Stop');
            };
            
            ws.onmessage = function(ev){
                var data = $.parseJSON(ev.data);
                $answers.queue(function(){ addAnswer(data); $(this).dequeue(); });
                $histogram.queue(function(){ updateHistogram(data); $(this).dequeue(); });
            };
            
            ws.onerror = function(){
                logDebug('An error has occurred');
                addAnswer('Something gone wrong...');
            };
            
            ws.onclose = function(ev){
                logDebug('Connection is closed (' + ev.code + ')');
                
                $button.attr('data-opened', 'false').removeAttr('disabled').val('Start');
            };
            
            /* TESTing (answer adding)
            var test = setInterval(function(){
                var rand = Math.random()*5;
                for(var i = 0; i < rand; i++){ addAnswer('test ' + Date.now()); }
            }, 100);
            setTimeout(function(){clearInterval(test)}, 30000);
            //*/
            
        } else {
            
            $button.attr('disabled','disabled').val('Disconnecting...');
            
            ws.close();
            
            $answers.removeClass('autoscroll');
            
        }
        
    }
    
    function filterN(){
        var defval = '1';
        var val = $input.val().replace(/^0+/, '').replace(/[^0-9]/gi, '');
        val = val == '' ? defval : val; 
        $input.val(val);
    }
    
    function randN(){
        $input.val(Math.logRand(n_minimum, n_maximum));
    }
    
    function example(){
        var value = $input.val();
        var value = value ? value : '108';
        var result = Math.collatzSequence(value);
        $example.html('Path for '+ value +' is ' + result.join(' &gt; '));
    }
    
    $button.attr('data-opened', 'false');
    
    $button.on('click',  function(ev){
        filterN();
        doAction();
    });

    $example.on('click',  function(ev){
        $example.toggleClass('collapsed');
    });
    
    $input.on('change', function(ev){
        filterN();
        example();
    });
    
    $input.on('keyup', function(ev){
        // Backspace
        if(ev.keyCode != 8){ filterN(); }
        // Enter
        if(ev.keyCode == 13){ doAction(); }
        example();
    });
    
    $random.on('click', function(ev){
        randN();
        example();
    });
    
    // randN();
    example();
    
});
