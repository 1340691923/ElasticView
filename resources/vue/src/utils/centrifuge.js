import  { Centrifuge } from "centrifuge";
import {getToken} from "@/utils/auth";

let centrifuge = null;

let subscriptions = new Map();

function getBaseUrl(){

  let href = window.location.href;
  if(!import.meta.env.PROD){
    href = import.meta.env.VITE_APP_API_URL
  }

  const protocol = href.split('//')[0] === 'http:' ? 'ws' : 'wss';
  const host = href.split('//')[1].split('/')[0];
  return `${protocol}://${host}/ws`
}

/**
 * 初始化 Centrifuge 实例（仅在用户登录后调用）
 */
export const initCentrifuge = () => {
  if (centrifuge) {
    //console.warn("Centrifuge is already initialized.");
    return;
  }

  if(getToken()==null || getToken() == "" ){
    return
  }

  centrifuge = new Centrifuge(getBaseUrl(), {
    token:getToken(),
    debug: true,
    timeout: 5000,
    websocket: null
  });

  // 连接状态处理
  centrifuge.on('connecting', function(ctx) {
    console.log('正在连接...', ctx);
  });

  centrifuge.on('connected', function(ctx) {
    console.log('已连接！', ctx);
  });

  centrifuge.on('disconnected', function(ctx) {
    console.log('连接断开', ctx);
  });

  centrifuge.on('error', function(ctx) {
    console.error('连接错误', ctx);
  });

  centrifuge.on('message', function (message) {

    let msgChannel = message.data.channel
    let arr = msgChannel.split("$v$");
    if(arr.length != 2) return

    subscriptions.forEach((cb, channel) => {
      if(channel == msgChannel){
        cb(message.data.data)
      } 
    });


  }); 

  centrifuge.connect();
}; 
 
/**
 * 订阅一个频道
 * @param {string} channel 频道名称
 * @param {function} callback 频道返回的消息
 */
export const SubscribeToChannel = (channel, publicationCb) => {
  if (!centrifuge) {
    console.error("长连接未建立.");
    return;
  }

  if (centrifuge.getSubscription(channel)) {
    console.warn(` 已经订阅该频道 ${channel}`);
    return;
  }

  // 订阅聊天室频道
  const sub = centrifuge.newSubscription(channel);

  sub.on('publication', function(ctx) {
    publicationCb(ctx.data)
  });

  sub.on('subscribing', function(ctx) {
    console.log('正在订阅频道...', ctx);
  });

  sub.on('subscribed', function(ctx) {
    console.log('已订阅频道！', ctx);
  });

  sub.on('unsubscribed', function(ctx) {
    console.log('取消订阅频道', ctx);
  });

  sub.on('error', function(ctx) {
    console.error('订阅错误', ctx);
  });

  // 订阅频道
  sub.subscribe();

  subscriptions.set(channel,publicationCb);
};

export const publish = (channel,data,successCb,errCb) => {

  if (!centrifuge) {
    console.error("长连接未初始化.");
    return;
  }

  if (!centrifuge.getSubscription(channel)) {
    console.warn(`请先订阅频道 ${channel}`);
    return;
  }

  let sub = centrifuge.getSubscription(channel)

  sub.publish(data).then(function(res) {
    successCb(res)
  }).catch(function(err) {
    errCb(err)
  });
}

/**
 * 取消订阅频道
 * @param {string} channel 频道名称
 */
export const unsubscribeFromChannel = (channel) => {
  if (!centrifuge || !centrifuge.getSubscription(channel)) return;
  const subscription = centrifuge.getSubscription(channel);
  if (subscription) {
    subscription.unsubscribe(() => {
      console.log('Unsubscribed successfully',channel);
    });
    centrifuge.removeSubscription(subscription)
  } else {
    console.log('No active subscription');
  }
  subscriptions.delete(channel);
  console.log(`停止订阅该频道 ${channel}`);
};

/**
 * 断开 Centrifuge 连接（用户退出时调用）
 */
export const disconnectCentrifuge = () => {
  if (centrifuge) {
    subscriptions.forEach((cb, channel) => unsubscribeFromChannel(channel));
    centrifuge.disconnect();
    centrifuge = null;
    console.log("断开长连接.");
  }
};
