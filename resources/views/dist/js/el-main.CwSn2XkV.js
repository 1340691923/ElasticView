import{d as e,c6 as t,b as s,c as a,e as r,h as n,k as i,i as o,g as l,_ as c,n as u,q as d,H as h}from"./index.DhFaMfT1.js";const f=e({name:"ElContainer"});var p=c(e({...f,props:{direction:{type:String}},setup(e){const c=e,u=t(),d=s("container"),h=a((()=>{if("vertical"===c.direction)return!0;if("horizontal"===c.direction)return!1;if(u&&u.default){return u.default().some((e=>{const t=e.type.name;return"ElHeader"===t||"ElFooter"===t}))}return!1}));return(e,t)=>(r(),n("section",{class:o([l(d).b(),l(d).is("vertical",l(h))])},[i(e.$slots,"default")],2))}}),[["__file","container.vue"]]);const g=e({name:"ElAside"});var m=c(e({...g,props:{width:{type:String,default:null}},setup(e){const t=e,c=s("aside"),d=a((()=>t.width?c.cssVarBlock({width:t.width}):{}));return(e,t)=>(r(),n("aside",{class:o(l(c).b()),style:u(l(d))},[i(e.$slots,"default")],6))}}),[["__file","aside.vue"]]);const v=e({name:"ElFooter"});var _=c(e({...v,props:{height:{type:String,default:null}},setup(e){const t=e,c=s("footer"),d=a((()=>t.height?c.cssVarBlock({height:t.height}):{}));return(e,t)=>(r(),n("footer",{class:o(l(c).b()),style:u(l(d))},[i(e.$slots,"default")],6))}}),[["__file","footer.vue"]]);const y=e({name:"ElHeader"});var E=c(e({...y,props:{height:{type:String,default:null}},setup(e){const t=e,c=s("header"),d=a((()=>t.height?c.cssVarBlock({height:t.height}):{}));return(e,t)=>(r(),n("header",{class:o(l(c).b()),style:u(l(d))},[i(e.$slots,"default")],6))}}),[["__file","header.vue"]]);const b=e({name:"ElMain"});var $=c(e({...b,setup(e){const t=s("main");return(e,s)=>(r(),n("main",{class:o(l(t).b())},[i(e.$slots,"default")],2))}}),[["__file","main.vue"]]);const k=d(p,{Aside:m,Footer:_,Header:E,Main:$});h(m),h(_),h(E),h($);export{k as E};