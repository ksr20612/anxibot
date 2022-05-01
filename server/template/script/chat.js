/* ------------ timeLineMax (animations) ------------- */
const detach_message = new TimelineMax();
const detach_motion_speed = 0.2;

const fill_message = new TimelineMax();
const fill_motion_speed = 0.2;

function detachMsg(){
    const $msg = $(".loading").last().parent().parent().parent();
    detach_message.to($msg, detach_motion_speed, {y : 10, autoAlpha : 0})
    .call(remove, [$msg]);
}

var remove = function removeLastMsg(target){
    target.remove();
}

/* ------------- static variables --------------- */
let msgIdxBot = 0;
let msgIdxUser = 0;
const botName = "챗봇";
var userName;
let nowLoading = false;
let textType = 0; // 고정 시나리오 순서 번호 : 0 ~ 5

/* --------- websocket initialization ----------- */
const scheme = window.location.protocol == "https:" ? 'wss://' : 'ws://';
const uri = scheme + window.location.hostname + (location.port ? ':' + location.port : '') + '/talk';
const websocket = new WebSocket(uri);
websocket.onopen = () => {
    console.log('Connected');
};
websocket.onclose = () => {
    console.log('Closed');
};
// it got msg
websocket.onmessage = (e) => {
    
    const result = JSON.parse(e.data);
    console.log(`answer : ${result.answer}`);
    console.log(`typeChange : ${result.typeChange}`);
    
    if(result.answer && result.typeChange){

        changeType(result.typeChange);
        setTimeout(()=>{
            appendMsg(1,result.answer);
        },1500);
    
    }else {
        alert("챗봇이 메세지를 보내지 못했습니다...\n다시한번 이야기해주세요.");
    }

};
websocket.onerror = (e) => {
    console.log('Error (see console)');
    console.log(e);
    alert("메시지 전송에 실패하였습니다.");
};

/* ------------- append messages --------------- */
// append messages with loading dots
function appendMsg(sender, sentence){

    let html='';
    if(sentence){
        nowLoading = false;
        const $msg = $(".loading").last().parent();
        $msg.html(sentence);
        $msg.css("min-width","initial");
        return true;
    }

    // sender : { "bot" : 1, "me" : 0 }
    if(sender){

        msgIdxBot += 1;
        html += '<div class="message_bot" id="message_bot_'+msgIdxBot+'">';
        html +=     '<div class="bot_img"></div>';
        html +=     '<div class="bot">';
        html +=         '<div class="bot_name">'+botName+'</div>';
        html +=         '<div class="bot_msg"><div class="loading"></div></div>';
        html +=         '</div>'
        html +=     '</div>';
        html += '</div>';

    }else {

        msgIdxUser += 1;
        html += '<div class="message_user" id="message_user_'+msgIdxUser+'">';
        html +=     '<div class="user">';
        if(userName){
            html +=         '<div class="user_name">'+userName+'</div>';
        }else {
            html +=         '<div class="user_name">???</div>';            
        }
        html +=         '<div class="user_msg"><div class="loading"></div></div>';
        html +=     '</div>';
        html +=     '<div class="user_img"></div>';
        html += '</div>';

    }

    $(".screen_box").append(html);
    $(".loading").last().parent().css("min-width","6rem");

    // # scroll down to the bottom
    $('.screen_box').scrollTop($('.screen_box')[0].scrollHeight);

}

function sendMsg(){

    const msg = $("textarea").val();

    // 첫 두 발화에 대한 응답이 0, 그 이후로 한 개씩 증가
    //const textType = (msgIdxUser===0||msgIdxUser===1)? 0 : msgIdxUser-2;
    
    if(msg.length !== 0){

        // # 1 append message
        appendMsg(0,msg);

        // # 2 get token & paramData
        var paramData = {
            "type" : textType,
            "idx" : msgIdxUser,
            "msg" : msg,
            "date" : getDate()
        };
        var token = localStorage.getItem("token");

        // # 3 send data to the server
        setTimeout(function(){
            appendMsg(1,"");
        },500);
        $("textarea").val("");
        console.log(paramData);
        websocket.send(paramData);
        // $.ajax({
        //     type : "POST",
        //     url : "/api/reply",
        //     contentType : "application/json;charset=utf-8",
        //     dataType : "json",
        //     cache : false,
        //     data :  JSON.stringify(paramData),
        //     beforeSend: function (xhr) {
        //         appendMsg(1,"");
        //         xhr.setRequestHeader("Authorization","JWT " + token);
        //     },
        //     success : function(result){
    
        //         $(".bot .bot_msg").last().html(result.reply);

        //     },
        //     error : function(result){
        //         alert("전송에 실패하였습니다.");
        //         $(".message_user").last().remove();
        //         $(".message_bot").last().remove();
        //     },
        //     complete : function(result) {

        //         $("textarea").val("");

        //     }
        // });  

    }
}

function delay(sec){
    return new Promise(resolve => {
        setTimeout(()=>{},sec);
    });
}

const getDate = () => {
    const today = new Date();   
    const year = today.getFullYear(); // 년도
    const month = (today.getMonth() + 1 < 10)? "0"+String(today.getMonth()+1) : ""+String(today.getMonth()+1); 
    const date = (today.getDate() < 10)? "0"+(today.getDate()) : ""+today.getDate(); 

    return year + "-" + month + "-" + date;
}

const changeType = (param) => {
    return textType + (param*1);
}

$(document).ready(function () {

    $("textarea").on('keyup keypress',function(e){
        var keyCode = e.keyCode || e.which;
        if(keyCode == 13){
            e.preventDefault();
            
            // submit
            sendMsg();

            return false;

        }
    })

    // if typing message, loading dot emerges
    $("textarea").on('propertychange change keyup paste input', (e) => {

        var currentVal = $("textarea").val();
        if(!nowLoading && currentVal.length != 0){
            nowLoading = true;
            console.log("off -> on : " + nowLoading);
            appendMsg(0,"");
        }else if(nowLoading && currentVal.length == 0){
            detachMsg();
            nowLoading = false;
            console.log("on -> off : " + nowLoading);
        }

    });

});
