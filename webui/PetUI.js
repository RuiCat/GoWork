/* jshint esversion: 6 */
console.log('%c PetUI %c 网页UI框架\n %c作者: RuiCat',
    'font-family: "Helvetica Neue",Helvetica, Arial, sans-serif;font-size:64px;color:#FFB6C1;-webkit-text-fill-color:#FFB6C1;-webkit-text-stroke: 1px #FFB6C1;',
    'font-size:25px;color:#FFB6C1;', 'font-size:15px;color:#FFB6C1;');
// 框架扩展
const PetUIExpand = (typeof this.PetUIExpand != "undefined") ? PetUIExpand : {
    // 打字效果
    Typing: function () {
        var Ty = Object(function (str, obj, callback, T) {
            var t = T ? T : 100;
            if (Ty.I <= str.length) {
                obj.innerText = str.slice(0, Ty.I++);
                var timer = window.setTimeout(function () {
                    Ty(str, obj, callback, t);
                    window.clearTimeout(timer);
                }, t); //递归调用
            } else {
                Ty.I = 0;
                if (callback == null) {
                    return;
                }
                if (typeof (callback) !== "undefined") {
                    callback();
                }
            }
        });
        Ty.I = 0;
        return Ty;
    },
    // 获得 文本/字节集 的 Src
    SrcFile: function (data) {
        var aFilePath = [];
        aFilePath.push(data);
        var blob = new Blob(aFilePath, {
            type: 'application/octet-binary'
        });
        var url = window.URL.createObjectURL(blob);
        url.Revoke = () => { window.URL.revokeObjectURL(url); };
        return url
    },
    // Blob 分割
    SliceBlob: function (blob, start, end, type) {
        var type = type || blob.type;
        if (blob.mozSlice) {
            return blob.mozSlice(start, end, type);
        } else if (blob.webkitSlice) {
            return blob.webkitSlice(start, end, type);
        } else {
            throw new Error("This doesn't work!");
        }
    },
    // 后台处理模块
    DaemonJS: function (JsData, CallBack) {
        var SrcFile = new PetUIExpand.SrcFile(`var $ = this,$Call = console.log; \n
        $.log = $.postMessage;
        $.onmessage = (e) => $Call(e.data,e); \n 
        $.call = (e) => $.postMessage({call:e}) \n` + JsData);
        if (typeof (Worker) !== "undefined") {
            var $ = new Worker(SrcFile.src)
            $.data = []
            $.CallBack = CallBack ? CallBack : (...e) => console.log(...e);
            $.SrcFile = SrcFile;
            // 返回数据处理
            $.onmessage = (e) => {
                var data = e.data;
                if (data.call) {
                    if (data.call == 'Destroy') {
                        $.Destroy();
                    }
                } else {
                    $.data.push(data);
                    if ($.CallBack) {
                        $.CallBack(data);
                    }
                }
            };
            // 传递消息
            $.Call = $.postMessage;
            // 错误处理
            $.onerror = (error) => {
                console.warn("脚本执行错误:", error);
                $.Destroy();
            };
            // 设置销毁
            $.Destroy = () => {
                $.terminate();
                $.data.length = 0;
            };
            return $;
        } else {
            console.log("无法创建后台任务");
        }
        return undefined;
    },
    // 单击转双击判定
    DoubleClick: function (Click, Double, Mousedown, Mouseup) {
        var mouseup = null;
        var $ = (...e) => {
            if (typeof (Mousedown) == "function") {
                Mousedown(...e);
            }
            if ($.look) {
                $.Dblclick = true;
                return;
            }
            $.look = true;
            var timer = window.setTimeout(function () {
                if (!($.look)) {
                    return;
                }
                if ($.Dblclick) {
                    if (typeof (Double) == "function") {
                        Double(...e);
                    }
                } else {
                    if (typeof (Click) == "function") {
                        Click(...e);
                    }
                }
                $.look = false;
                $.Dblclick = false;
                if (mouseup != null && typeof (Mouseup) == "function") {
                    Mouseup(...mouseup);
                    mouseup = null;
                }
                window.clearTimeout(timer);
            }, 250);
        };
        $.Mouseup = (...e) => {
            if ($.look && $.Dblclick) {
                if (typeof (Mouseup) == "function") {
                    Mouseup(...e);
                }
            } else {
                mouseup = e;
            }
        };
        $.look = false;
        $.Dblclick = false;
        return $;
    },
    // 不规则矩形判定
    IsPointInPolygon: function (x, y, coords, width, height) {
        var wn = 0;
        width = width ? width : 1;
        height = height ? height : 1;
        for (var shiftP, shift = coords[1] * height > y, i = 3; i < coords.length; i += 2) {
            shiftP = shift;
            shift = coords[i] * height > y;
            if (shiftP != shift) {
                var n = (shiftP ? 1 : 0) - (shift ? 1 : 0);
                // dot product for vectors (c[0]-x, c[1]-y) . (c[2]-x, c[3]-y)
                if (n * ((coords[i - 3] * width - x) * (coords[i - 0] * height - y) - (coords[i - 2] * height - y) * (coords[i - 1] * width - x)) > 0) {
                    wn += n;
                }
            }
        }
        return wn;
    },
    // 标记获取
    Guid: function () {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            var r = Math.random() * 16 | 0,
                v = c == 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    },
    // 平台判定 
    IsPC: function () {
        var userAgentInfo = navigator.userAgent;
        var Agents = ["Android", "iPhone", "SymbianOS", "Windows Phone", "iPad", "iPod"];
        var flag = true;
        for (var v = 0; v < Agents.length; v++) {
            if (userAgentInfo.indexOf(Agents[v]) > 0) {
                flag = false;
                break;
            }
        }
        return flag;
    },
    // 是否支持触摸
    HasTouch: function () {
        var touchObj = {};
        touchObj.isSupportTouch = "ontouchend" in document ? true : false;
        touchObj.isEvent = touchObj.isSupportTouch ? true : false;
        return touchObj.isEvent;
    },
    // 使对象事件重复定义触发功能
    SetChangeEvent: function (obj, event) {
        var $ = {};
        $.Proto = [];
        $.TheDefaultReturn = null;
        $.Transmit = false;
        $.Destroy = function (id) {
            for (var i in $.Proto) {
                if ($.Proto[i] == id) {
                    $.Proto.splice(i, 1);
                    return true;
                }
            }
            return false;
        };
        // 保存原函数指针
        try {
            if (typeof obj[event] == "function") {
                $.Proto.push(obj[event]);
            }
            var ret = (...name) => {
                var Ret, i;
                if ($.TheDefaultReturn != undefined) {
                    Ret = $.Proto[$.TheDefaultReturn].call(obj, ...name);
                    return Ret;
                }
                try {
                    $.Proto.forEach(element => {
                        element.call(obj, ...name);
                        // 传递返回
                        if ($.Transmit) {
                            $.Transmit = false;
                            return;
                        }
                    });
                } catch (err) {
                    console.warn(err, "\nID: " + i);
                    $.Proto.splice(i, 1);
                }
            };
            ret.__proto__ = $;
            Object.defineProperty(obj, event, {
                // 改变事件 添加到事件触发数组
                set: function (name) {
                    if (typeof name == "function") {
                        obj[event + "ID"] = $.Proto.length;
                        $.Proto.push(name);
                    }
                },
                // 调用事件 循环调用事件
                get: function () {
                    if ($.Proto == 1) {
                        return $.Proto[0];
                    }
                    // 返回 构建数组
                    return ret;
                }
            });
        } catch (err) {
            console.warn(err);
            return null;
        }
        return $;
    },
    // 对象添加初始化
    SetOnload: function (obj) {
        var $ = function (...n) {
            $.Param = n;
            $.lock = true;
            for (var i in $.Proto) {
                $.Proto[i](...n);
            }
            $.Proto.splice(0, $.Proto.length);
        };
        $.__proto__ = {
            lock: false,
            Proto: [],
            Param: [],
        };
        // 保存之前需要初始化的函数
        if (typeof obj.onload == "function") {
            $.Proto.push(obj.onload);
        }
        // 设置初始化
        obj["onload"] = $;
        // 拦截之后初始化函数设置
        Object.defineProperty(obj, 'onload', {
            // 改变事件 添加到事件触发数组
            set: function (fn) {
                if (typeof fn != "function") {
                    return;
                }
                if ($.lock) {
                    fn(...$.Param);
                } else {
                    $.Proto.push(fn);
                }
            },
            get: function () {
                console.warn("函数: onload -> 错误调用 重复初始化");
            }
        });
    },
    /**
     * HttpQuest 通信
     * @param option     连接参数 { url , method , data , timeout , responseType }
     * @param callback   回调函数 (err, result) => { }
     */
    HttpQuest: (option, callback = (err, result) => { }) => {
        var url = option.url;
        var method = option.method;
        var data = option.data;
        var timeout = option.timeout || 0;
        var xhr = null;
        try {
            xhr = new ActiveXObject("Msxml2.XMLHTTP");
        } catch (e) {
            try {
                xhr = new ActiveXObject("Microsoft.XMLHTTP");
            } catch (e1) {
                xhr = new XMLHttpRequest();
            }
        }
        // 返回数据类型
        xhr.responseType = (([
            'text',
            'arraybuffer',
            'blob',
            'document'
        ]).indexOf(option.responseType) != -1) ? option.responseType : 'text';
        if (timeout > 0) {
            xhr.timeout = timeout;
        }
        // 设置回调函数
        xhr.onreadystatechange = () => {
            if (xhr.readyState == 4) {
                if (xhr.status >= 200 && xhr.status < 400) {
                    var result = '';
                    switch (xhr.responseType) {
                        case 'text':
                            result = xhr.responseText;
                            try {
                                result = JSON.parse(xhr.responseText);
                            } catch (e) { }
                            break;
                        case 'document':
                            result = xhr.responseXML;
                            break;
                        default:
                            result = xhr.response;
                    }
                    if (callback) {
                        callback(null, result);
                    }
                } else {
                    if (callback) {
                        callback('status: ' + xhr.status);
                    }
                }
            }
        };
        // 打开连接
        xhr.open(method, url, true);
        if (typeof data === 'object') {
            try {
                data = JSON.stringify(data);
            } catch (e) { };
        }
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;");
        // 发送数据
        xhr.send(data);
        xhr.ontimeout = function () {
            if (callback) {
                callback('timeout');
            }
            console.log('%c连%c接%c超%c时', 'color:red', 'color:orange', 'color:purple', 'color:green');
        };
    },
    /**
     * HttpGet 通信
     * @param url        连接地址
     * @param callback   回调函数 (err, result) => { }
     */
    HttpGet: (url, callback) => {
        var option = url.url ? url : {
            url: url
        };
        option.responseType = (([
            'text',
            'arraybuffer',
            'blob',
            'document'
        ]).indexOf(url.responseType ? url.responseType : "") != -1) ? url.responseType : 'text';
        option.method = 'get';
        PetUIExpand.HttpQuest(option, callback);
    },
    /**
     * HttpPost 通信
     * @param option     连接参数 { url , method , data , timeout , responseType }
     * @param callback   回调函数 (err, result) => { }
     */
    HttpPost: (option, callback) => {
        option.method = 'post';
        PetUIExpand.HttpQuest(option, callback);
    },
    /**
     * 前端绑定API实现
     * @param callFn    注册本地函数列表
     * @param config    配置信息 onopen:连接回掉,url:连接地址
     * @returns         构建的后端结构体
     */
    Bind: function (callFn = {}, config = { onopen: undefined, url: "" }) {
        var obj = { url: "", value: {} };
        // 获取当前脚本的链接地址
        if (config?.url ?? "" != "") {
            obj.url = config.url
        } else if (document.scripts.length > 0) {
            obj.url = document.scripts[document.scripts.length - 1].src;
            obj.url = obj.url.substring(0, obj.url.lastIndexOf("."));
            var i = obj.url.indexOf("://");
            if (i > 0) { obj.url = obj.url.substring(i + 3); };
        } else {
            console.log("获取 API 调用接口失败");
            return undefined;
        };
        obj.ws = new WebSocket("ws://" + obj.url);
        obj.Send = function (type, name, value) {
            obj.ws.send(JSON.stringify({ type: String(type), name: String(name), value: value }));
        };
        obj.ws.onerror = obj.ws.onclose = function (etv) { obj = undefined; };
        obj.Close = function () { obj.ws.close(); };
        obj.ws.onmessage = function (e) {
            var eve = JSON.parse(e.data);
            switch (eve.type) {
                case "CallFunc":
                    callFn[eve.name](eve.value);
                    break;
                case "SetValue":
                    obj.value[eve.name] = eve.value;
                    break;
                case "DeleteValue":
                    delete obj.value[eve.name];
                    break;
                case "SetFunc":
                    obj.value[eve.name] = (...argumentsList) => obj.Send("CallFunc", eve.name, argumentsList);
                    break;
                case "DeleteFunc":
                    delete obj.value[eve.name];
                    break;
                case "onopen":
                    if (config.onopen != undefined) { config.onopen() };
                    break;
                default:
                    break;
            }
        };
        return new Proxy(obj.value, {
            getPrototypeOf(_) { return obj; },
            deleteProperty: function (_, prop) {
                if (typeof obj.value[prop] == "function") {
                    obj.Send("DeleteFunc", prop);
                } else {
                    obj.Send("DeleteValue", prop);
                }
                delete obj.value[prop];
            },
            set: function (_, prop, value) {
                if (typeof obj.value[prop] != "function") {
                    obj.value[prop] = value;
                    obj.Send("SetValue", prop, value);
                }
            },
            get: function (_, prop) {
                return obj.value[prop];
            },
        });
    },
    // 禁止菜单显示
    BanMenu: function () {
        document.oncontextmenu = function (evt) {
            var event = evt || window.event;
            if (event && event.returnValue) {
                event.preventDefault();
            } else {
                event.returnValue = false;
            }
        };
    },
};
// 主框架
const PetUI = (typeof this.PetUI != "undefined") ? this.PetUI : (() => {
    // 移动端事件BUG修复代码
    PetUIExpand.SetOnload(window);
    (() => {
        var touch = (event) => {
            var touch = [];
            event = event || window.event;
            touch = (event.touches && event.changedTouches)[0];
            touch.buttons = 1;
            touch.event = event;
            var $ = touch.target.Element;
            if ($) {
                if ($.TheMobileEventTakesEffect) {
                    switch (event.type) {
                        case "touchstart":
                            $.MobileEndEventPassing.Click(touch);
                            break;
                        case "touchend":
                            $.MobileEndEventPassing.Mouseup(touch);
                            break;
                        case "touchmove":
                            // 判断默认行为是否可以被禁用
                            if (event.cancelable) {
                                // 判断默认行为是否已经被禁用
                                if (!event.defaultPrevented) {
                                    event.preventDefault();
                                }
                            }
                            $.MobileEndEventPassing.Mousemove(touch);
                            break;
                    }
                }
            }
        };
        // 事件注册
        if (document.addEventListener) {
            document.addEventListener('touchstart', touch, false);
            document.addEventListener('touchmove', touch, false);
            document.addEventListener('touchend', touch, false);
        } else if (document.attachEvent) {
            document.attachEvent('ontouchstart', touch, false);
            document.attachEvent('ontouchmove', touch, false);
            document.attachEvent('ontouchend', touch, false);
        } else {
            document.ontouchstart = touch;
            document.ontouchmove = touch;
            document.ontouchend = touch;
        }
    })();
    // 继承函数
    var $ = {
        // UI - Style 事件拦截
        UI_Tackl_Style: function (style, Call = { set: () => { }, get: () => { } }) {
            // 触发事件回调
            var $ = this;
            var st = style;
            Object.defineProperty($.style, st, {
                set: function (name) {
                    $.style = $.style.cssText + st + ":" + name;
                    Call.set();
                },
                get: function () {
                    var style = $.style.cssText;
                    var arr = style.match(RegExp(st + ":(.*?);"));
                    if (arr != null) {
                        arr = arr[0].match(/:\s(.*?);/);
                        return arr[arr.length - 1].toString();
                    }
                    Call.get();
                    return "";
                }
            });
            return Call;
        },
        // style 设置函数
        CSS: function (style) {
            var typ = typeof style;
            if (typ == "object") {
                var css = "";
                if (this.style) {
                    for (var x in style) {
                        css += ";" + x + ":" + style[x];
                    }
                }
                if (css != "") {
                    this.style.cssText += css;
                }
            } else if (typ == "string") {
                this.style = this.cssText + ";" + style;
            }
            return this;
        },
        // 删除元素函数
        Destroy: function () {
            // 处理元素删除
            (this.parentNode).removeChild(this);
        },
        // 添加到元素
        Child: function (aChild, state) {
            if (aChild == undefined) {
                if (state && document.body.childNodes[0] != undefined) {
                    document.body.insertBefore(this, document.body.childNodes[0]);
                } else {
                    document.body.appendChild(this);
                }
                return this;
            }
            if (this.appendChild) {
                if (state && this.childNodes[0] != undefined) {
                    this.insertBefore(aChild, this.childNodes[0]);
                } else {
                    this.appendChild(aChild);
                }
            }
            return this;
        },
        // 样式
        Class: function (style, state) {
            if (typeof style == "string") {
                this.classList.toggle(style, state);
            } else {
                for (var i = 0; i < style.length; i++) {
                    this.classList.toggle(style[i], state);
                }
            }
        },
        // 事件注册
        AddHandler: function (type, handler, ...use) {
            if (this.addEventListener) {
                this.addEventListener(type, handler, use ? use : false);
            } else if (this.attachEvent) {
                this.attachEvent("on" + type, handler, use ? use : false);
            } else {
                this["on" + type] = handler;
            }
        },
        // 元素隐藏/显示
        Display: function (state) {
            if (state == undefined) {
                return this.style.display != 'none';
            }
            if (state) {
                this.style.display = 'inline';
            } else {
                this.style.display = 'none';
            }
            return state;
        },
        // 保证坐标
        GuaranteeCoords: function () {
            var Gua = { Offset: { Top: 0, Left: 0 } };
            var Top = 0,
                Left = 0,
                ParentNode = null;
            this.GuaranteeCoords = Gua;
            // 开启还原
            Gua.EnablementUndo = false;
            // 设置
            Gua.Save = () => {
                if (Gua.EnablementUndo) {
                    ParentNode = this.UI_Ele_Father && this.UI_Ele_Father != document.body ? this.UI_Ele_Father : {
                        offsetHeight: this.Position.Win.Height, // 高度
                        offsetWidth: this.Position.Win.Width // 宽度
                    };
                    Top = (this.offsetTop + parseInt(this.offsetHeight / 2)) / ParentNode.offsetHeight;
                    Left = (this.offsetLeft + parseInt(this.offsetWidth / 2)) / ParentNode.offsetWidth;
                }
            };
            // 保证记录坐标
            (this.UI_Tackl_Style("top")).set = Gua.Save;
            (this.UI_Tackl_Style("left")).set = Gua.Save;
            // 还原
            Gua.Recovery = () => {
                if (Gua.EnablementUndo) {
                    ParentNode = this.UI_Ele_Father && Ele.UI_Ele_Father != document.body ? this.UI_Ele_Father : {
                        offsetHeight: this.Position.Win.Height,
                        offsetWidth: this.Position.Win.Width
                    };
                    this.style = this.style.cssText + 'top:' + parseInt(Top * ParentNode.offsetHeight - parseInt(this.offsetWidth / 2) + Gua.Offset.Top) + "px;left:" + parseInt(Left * ParentNode.offsetWidth - parseInt(this.offsetHeight / 2) + Gua.Offset.Left) + "px;";
                }
            };
            // 自适应
            this.Position.Affair = Gua.Recovery;
            // 销毁自适应
            Gua.RemoveAdaptio = () => (this.Position.AddAffair.Destroy(this.Recovery));
            return Gua;
        },
        // 开启移动
        MobileEvent: function () {
            // 注册事件
            this.Set_Element();
            var Mobile = {
                // 移动锁
                MouseDown: false,
                // 是否可移动
                EnablementRemovable: false,
                // 移动回调
                TheMouseCallback: () => { },
            };
            // 移动坐标
            var Page = {
                X: 0,
                Y: 0,
            };
            // 设置事件
            this.Mousedown = function (e) {
                if (Mobile.MouseDown) {
                    return;
                }
                if (Mobile.EnablementRemovable) {
                    // 将Mouse事件锁定到指定元素上
                    e.setCapture && e.setCapture();
                    Mobile.MouseDown = true;
                    // 记录坐标
                    Page.X = e.clientX - parseInt(this.offsetLeft);
                    Page.Y = e.clientY - parseInt(this.offsetTop);
                }
            };
            // 释放事件
            this.Mouseup = function (e) {
                Mobile.MouseDown = false;
                this.GuaranteeCoords.Save && this.GuaranteeCoords.Save();
                e.releaseCapture && e.releaseCapture();
                this.Mousemove.TheDefaultReturn = null;
            };
            // 移动事件
            this.Mousemove = function (e) {
                e.preventDefault && e.preventDefault();
                var button = e.buttons;
                if (!Mobile.MouseDown) {
                    return;
                };
                if (button == 1 || button == 3) {
                    if (Mobile.EnablementRemovable) {
                        this.ElementInherit.Mousemove.Transmit = true;
                        var ParentNode = this.parentNode && this.parentNode != document.body ? this.parentNode : {
                            offsetWidth: this.Position.Win.Width,
                            offsetHeight: this.Position.Win.Height,
                            offsetTop: this.Position.WinSkewing.top,
                            offsetLeft: this.Position.WinSkewing.left
                        };
                        var X = (e.clientX - Page.X);
                        var Y = (e.clientY - Page.Y);
                        if (X + this.offsetWidth > ParentNode.offsetWidth ||
                            Y > ParentNode.offsetHeight - this.offsetHeight ||
                            Y < ParentNode.offsetTop ||
                            X < ParentNode.offsetLeft) {
                            return;
                        };
                        this.CSS({
                            top: Y + 'px',
                            left: X + 'px',
                        });
                        Mobile.TheMouseCallback && Mobile.TheMouseCallback(e);
                    }
                    return;
                };
                this.GuaranteeCoords.Save && this.GuaranteeCoords.Save();
                // 释放
                Mobile.MouseDown = false;
                e.releaseCapture && e.releaseCapture();
                this.Mousemove.TheDefaultReturn = null;
            };
            // 注册移动操作
            if (PetUIExpand.IsPC()) {
                document.Mousemove = this.Mousemove;
            }
            this.MobileEvent = Mobile;
            return Mobile;
        },
        // 创建元素
        Create: function (elementId, create) {
            // 创建对象
            var $ele = {};
            if (create) {
                $ele = document.getElementById(elementId);
            } else {
                if (typeof elementId == "string") {
                    $ele = document.createElement(elementId);
                } else {
                    if (typeof elementId.Create != "undefined") {
                        return elementId;
                    }
                    $ele = elementId;
                }
            }
            var ele = Object.assign($ele, $);
            // 支持事件类型
            ele.Click = null;     // 单击事件
            ele.Dblclick = null;  // 双击事件
            ele.Mousedown = null; // 鼠标按下
            ele.Mouseup = null;   // 释放事件
            ele.Mouseout = null;  // 移出事件(电脑端)
            ele.Mousemove = null; // 移动事件
            // 开启继承
            ele.ElementInherit = {
                Click: new PetUIExpand.SetChangeEvent(ele, "Click"),
                Dblclick: new PetUIExpand.SetChangeEvent(ele, "Dblclick"),
                Mousedown: new PetUIExpand.SetChangeEvent(ele, "Mousedown"),
                Mouseup: new PetUIExpand.SetChangeEvent(ele, "Mouseup"),
                Mouseout: new PetUIExpand.SetChangeEvent(ele, "Mouseout"),
                Mousemove: new PetUIExpand.SetChangeEvent(ele, "Mousemove"),
            };
            // 鼠标事件
            let click = PetUIExpand.DoubleClick(
                ele.Click,
                ele.Dblclick,
                ele.Mousedown,
                ele.Mouseup,
            );
            // 默认事件处理
            if (PetUIExpand.HasTouch()) {
                //  移动端事件传递
                ele.MobileEndEventPassing = {
                    Click: click,
                    Mouseup: click.Mouseup,
                    Mousemove: ele.Mousemove,
                };
                // 事件生效标志
                ele.TheMobileEventTakesEffect = false;
                ele.Set_Element = () => {
                    // 设置触摸事件
                    ele.TheMobileEventTakesEffect = true;
                };
                ele.Remove_Element = () => {
                    // 删除触摸事件
                    ele.TheMobileEventTakesEffect = false;
                };
            } else {
                ele.Set_Element = () => {
                    ele.AddHandler("click", click);
                    ele.AddHandler("mouseup", click.Mouseup);
                    ele.AddHandler("mouseout", ele.Mouseout);
                    ele.AddHandler("mousedown", ele.Mousedown);
                    ele.AddHandler("mousemove", ele.Mousemove);
                };
                ele.Remove_Element = () => {
                    ele.AddHandler("click", undefined);
                    ele.AddHandler("mouseup", undefined);
                    ele.AddHandler("mouseout", undefined);
                    ele.AddHandler("mousedown", undefined);
                    ele.AddHandler("mousemove", undefined);
                };
            };
            return ele;
        },
    };
    // 获得网页窗口坐标信息
    $.Position = (function () {
        var $ = {};
        $.Affair = undefined;
        $.Scroll = {
            X: 0,
            Y: 0,
        };
        $.Win = {
            Width: 0,
            Height: 0,
        };
        // 坐标偏移/限制
        $.WinSkewing = {
            top: 0,
            left: 0,
            Width: 0,
            Height: 0,
        };
        $.Uinc = function () {
            // 滚动条位置
            if (self.scrollY) {
                $.Scroll.Y = self.scrollY;
                $.Scroll.X = self.scrollX;
            } else if (document.documentElement && document.documentElement.scrollTop) {
                $.Scroll.Y = document.documentElement.scrollTop;
                $.Scroll.X = document.documentElement.scrollLeft;
            } else if (document.body) {
                $.Scroll.Y = document.body.scrollTop;
                $.Scroll.X = document.body.scrollLeft;
            }
            // 获取窗口宽度
            if (window.innerWidth)
                $.Win.Width = window.innerWidth;
            else if ((document.body) && (document.body.clientWidth)) {
                $.Win.Width = document.body.clientWidth;
            }
            // 获取窗口高度
            if (window.innerHeight)
                $.Win.Height = window.innerHeight;
            else if ((document.body) && (document.body.clientHeight)) {
                $.Win.Height = document.body.clientHeight;
            }
            // 通过深入 Document 内部对 body 进行检测，获取窗口大小
            if (document.documentElement && document.documentElement.clientHeight && document.documentElement.clientWidth) {
                $.Win.Height = document.documentElement.clientHeight;
                $.Win.Width = document.documentElement.clientWidth;
            }
            $.Win.Height = $.Win.Height + $.WinSkewing.Height;
            $.Win.Width = $.Win.Width + $.WinSkewing.Width;
            if (typeof $.Affair == "function") {
                $.Affair();
            }
        };
        // 设定更变条件
        window.onload = $.Uinc;
        window.addEventListener("scroll", $.Uinc);
        window.addEventListener("resize", $.Uinc);
        $.AddAffair = new PetUIExpand.SetChangeEvent($, 'Affair');
        return $;
    })();
    // 对 document.Mousemove 初始化
    $.Mousemove = (function () {
        var ele = { Mousemove: null };
        var event = function (e) {
            e = e || window.event;
            e.buttons = 1;
            if (e.changedTouches) {
                e.clientX = e.changedTouches[0].pageX;
                e.clientY = e.changedTouches[0].pageY;
            }
            return e;
        };
        if (PetUIExpand.IsPC()) {
            $.AddHandler.call(document, "mousemove", (e) => ele.Mousemove(e))
        } else {
            $.AddHandler.call(document, "touchmove", (e) => ele.Mousemove(event(e)), {
                passive: true
            });
        }
        ele.Mousemove = new PetUIExpand.SetChangeEvent(ele, 'Mousemove');
        return ele.Mousemove;
    })();
    var Create = function (...Param) {
        if (Param.length == 0) {
            return (id) => {
                switch (typeof Param[0]) {
                    case "string":
                        if (Param[0][0] == "#") {
                            return $.Create(Param[0].slice(1, Param[0].length), true);
                        } else {
                            return $.Create(Param[0]);
                        }
                    case "object":
                        return Create.call(this, id);
                    default:
                        return undefined;
                }
            }
        }
        var ele = this.Child ? this : undefined;
        var $ele = ele;
        Param.forEach(element => {
            if (typeof element == "string") {
                if (element[0] == "#") {
                    $ele = $.Create(element.slice(1, element.length), true);
                } else {
                    $ele = $.Create(element);
                };
                ele = ele ? ele.Child($ele) : ele = $ele;
            } else if (typeof element.nodeName != "undefined") {
                $ele = $.Create(element);
                ele = ele ? ele.Child($ele) : ele = $ele;
            } else if (element instanceof Array) {
                if ($ele == undefined) {
                    $ele = {
                        Create: (param) => {
                            ele = Create(param);
                            return ele;
                        }
                    };
                };
                var node = $ele;
                var $node = undefined;
                element.forEach(element => {
                    if (element instanceof Object) {
                        $node = Create.call($node, element);
                        node = $node;
                        return;
                    } else {
                        $node = Create(element);
                    }
                    if (ele == undefined) {
                        ele = $node;
                    } else {
                        node.Child($node)
                    }
                    node = $node;
                });
            } else if (element instanceof Object) {
                $ele = $ele ? $ele : ele;
                for (let k in element) {
                    if (element instanceof Function) {
                        $ele[k](...element[k]);
                    } else {
                        $ele[k] = element[k];
                    }
                }
            }
        });
        if (ele == undefined) {
            return this;
        }
        ele.Create = Create;
        return ele;
    };
    // 主框架实现
    return new Proxy(function (...Param) {
        try {
            return Create(...Param);
        } catch (err) {
            console.warn("PetUI\n错误类型:", err);
            return null;
        }
    }, {
        getPrototypeOf(_) { return $; },
        set: function (_, prop, value) { },
        get: function (_, prop) { return $[prop]; },
    });
})();