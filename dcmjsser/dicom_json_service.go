package main

import "io/ioutil"
import "net/http"
import "strconv"
import "log"
import "encoding/json"
import "errors"
import "time"
import "os"

import "encoding/base64"

const htmlData = "PCFkb2N0eXBlIGh0bWw+DQo8aHRtbD4NCg0KPGhlYWQ+DQoJPHRpdGxlPmR0b29scyBVSTwvdGl0bGU+DQoJPG1ldGEgbmFtZT0idmlld3BvcnQiIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCI+DQoJPGxpbmsgcmVsPSJzdHlsZXNoZWV0IiBocmVmPSJodHRwczovL25ldGRuYS5ib290c3RyYXBjZG4uY29tL2Jvb3Rzd2F0Y2gvMy4wLjAvc2xhdGUvYm9vdHN0cmFwLm1pbi5jc3MiPg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0IiBzcmM9Imh0dHBzOi8vYWpheC5nb29nbGVhcGlzLmNvbS9hamF4L2xpYnMvanF1ZXJ5LzIuMC4zL2pxdWVyeS5taW4uanMiPjwvc2NyaXB0Pg0KCTxzY3JpcHQgdHlwZT0idGV4dC9qYXZhc2NyaXB0IiBzcmM9Imh0dHBzOi8vbmV0ZG5hLmJvb3RzdHJhcGNkbi5jb20vYm9vdHN0cmFwLzMuMS4xL2pzL2Jvb3RzdHJhcC5taW4uanMiPjwvc2NyaXB0Pg0KCTxzdHlsZSB0eXBlPSJ0ZXh0L2NzcyI+DQoJCWJvZHkgew0KCQkJcGFkZGluZy10b3A6IDIwcHg7DQoJCX0NCgkJDQoJCS5mb290ZXIgew0KCQkJYm9yZGVyLXRvcDogMXB4IHNvbGlkICNlZWU7DQoJCQltYXJnaW4tdG9wOiA0MHB4Ow0KCQkJcGFkZGluZy10b3A6IDQwcHg7DQoJCQlwYWRkaW5nLWJvdHRvbTogNDBweDsNCgkJfQ0KCQkvKiBNYWluIG1hcmtldGluZyBtZXNzYWdlIGFuZCBzaWduIHVwIGJ1dHRvbiAqLw0KCQkNCgkJLmp1bWJvdHJvbiB7DQoJCQl0ZXh0LWFsaWduOiBjZW50ZXI7DQoJCQliYWNrZ3JvdW5kLWNvbG9yOiB0cmFuc3BhcmVudDsNCgkJfQ0KCQkNCgkJLmp1bWJvdHJvbiAuYnRuIHsNCgkJCWZvbnQtc2l6ZTogMjFweDsNCgkJCXBhZGRpbmc6IDE0cHggMjRweDsNCgkJfQ0KCQkvKiBDdXN0b21pemUgdGhlIG5hdi1qdXN0aWZpZWQgbGlua3MgdG8gYmUgZmlsbCB0aGUgZW50aXJlIHNwYWNlIG9mIHRoZSAubmF2YmFyICovDQoJCQ0KCQkubmF2LWp1c3RpZmllZCB7DQoJCQliYWNrZ3JvdW5kLWNvbG9yOiAjZWVlOw0KCQkJYm9yZGVyLXJhZGl1czogNXB4Ow0KCQkJYm9yZGVyOiAxcHggc29saWQgI2NjYzsNCgkJfQ0KCQkNCgkJLm5hdi1qdXN0aWZpZWQgPiBsaSA+IGEgew0KCQkJcGFkZGluZy10b3A6IDE1cHg7DQoJCQlwYWRkaW5nLWJvdHRvbTogMTVweDsNCgkJCWNvbG9yOiAjNzc3Ow0KCQkJZm9udC13ZWlnaHQ6IGJvbGQ7DQoJCQl0ZXh0LWFsaWduOiBjZW50ZXI7DQoJCQlib3JkZXItYm90dG9tOiAxcHggc29saWQgI2Q1ZDVkNTsNCgkJCWJhY2tncm91bmQtY29sb3I6ICNlNWU1ZTU7DQoJCQkvKiBPbGQgYnJvd3NlcnMgKi8NCgkJCQ0KCQkJYmFja2dyb3VuZC1yZXBlYXQ6IHJlcGVhdC14Ow0KCQkJLyogUmVwZWF0IHRoZSBncmFkaWVudCAqLw0KCQkJDQoJCQliYWNrZ3JvdW5kLWltYWdlOiAtbW96LWxpbmVhci1ncmFkaWVudCh0b3AsICNmNWY1ZjUgMCUsICNlNWU1ZTUgMTAwJSk7DQoJCQkvKiBGRjMuNisgKi8NCgkJCQ0KCQkJYmFja2dyb3VuZC1pbWFnZTogLXdlYmtpdC1ncmFkaWVudChsaW5lYXIsIGxlZnQgdG9wLCBsZWZ0IGJvdHRvbSwgY29sb3Itc3RvcCgwJSwgI2Y1ZjVmNSksIGNvbG9yLXN0b3AoMTAwJSwgI2U1ZTVlNSkpOw0KCQkJLyogQ2hyb21lLFNhZmFyaTQrICovDQoJCQkNCgkJCWJhY2tncm91bmQtaW1hZ2U6IC13ZWJraXQtbGluZWFyLWdyYWRpZW50KHRvcCwgI2Y1ZjVmNSAwJSwgI2U1ZTVlNSAxMDAlKTsNCgkJCS8qIENocm9tZSAxMCssU2FmYXJpIDUuMSsgKi8NCgkJCQ0KCQkJYmFja2dyb3VuZC1pbWFnZTogLW1zLWxpbmVhci1ncmFkaWVudCh0b3AsICNmNWY1ZjUgMCUsICNlNWU1ZTUgMTAwJSk7DQoJCQkvKiBJRTEwKyAqLw0KCQkJDQoJCQliYWNrZ3JvdW5kLWltYWdlOiAtby1saW5lYXItZ3JhZGllbnQodG9wLCAjZjVmNWY1IDAlLCAjZTVlNWU1IDEwMCUpOw0KCQkJLyogT3BlcmEgMTEuMTArICovDQoJCQkNCgkJCWZpbHRlcjogcHJvZ2lkOiBEWEltYWdlVHJhbnNmb3JtLk1pY3Jvc29mdC5ncmFkaWVudChzdGFydENvbG9yc3RyPScjZjVmNWY1JywgZW5kQ29sb3JzdHI9JyNlNWU1ZTUnLCBHcmFkaWVudFR5cGU9MCk7DQoJCQkvKiBJRTYtOSAqLw0KCQkJDQoJCQliYWNrZ3JvdW5kLWltYWdlOiBsaW5lYXItZ3JhZGllbnQodG9wLCAjZjVmNWY1IDAlLCAjZTVlNWU1IDEwMCUpOw0KCQkJLyogVzNDICovDQoJCX0NCgkJDQoJCS5uYXYtanVzdGlmaWVkID4gLmFjdGl2ZSA+IGEsDQoJCS5uYXYtanVzdGlmaWVkID4gLmFjdGl2ZSA+IGE6aG92ZXIsDQoJCS5uYXYtanVzdGlmaWVkID4gLmFjdGl2ZSA+IGE6Zm9jdXMgew0KCQkJYmFja2dyb3VuZC1jb2xvcjogI2RkZDsNCgkJCWJhY2tncm91bmQtaW1hZ2U6IG5vbmU7DQoJCQlib3gtc2hhZG93OiBpbnNldCAwIDNweCA3cHggcmdiYSgwLCAwLCAwLCAuMTUpOw0KCQl9DQoJCQ0KCQkubmF2LWp1c3RpZmllZCA+IGxpOmZpcnN0LWNoaWxkID4gYSB7DQoJCQlib3JkZXItcmFkaXVzOiA1cHggNXB4IDAgMDsNCgkJfQ0KCQkNCgkJLm5hdi1qdXN0aWZpZWQgPiBsaTpsYXN0LWNoaWxkID4gYSB7DQoJCQlib3JkZXItYm90dG9tOiAwOw0KCQkJYm9yZGVyLXJhZGl1czogMCAwIDVweCA1cHg7DQoJCX0NCgkJDQoJCUBtZWRpYShtaW4td2lkdGg6IDc2OHB4KSB7DQoJCQkubmF2LWp1c3RpZmllZCB7DQoJCQkJbWF4LWhlaWdodDogNTJweDsNCgkJCX0NCgkJCS5uYXYtanVzdGlmaWVkID4gbGkgPiBhIHsNCgkJCQlib3JkZXItbGVmdDogMXB4IHNvbGlkICNmZmY7DQoJCQkJYm9yZGVyLXJpZ2h0OiAxcHggc29saWQgI2Q1ZDVkNTsNCgkJCX0NCgkJCS5uYXYtanVzdGlmaWVkID4gbGk6Zmlyc3QtY2hpbGQgPiBhIHsNCgkJCQlib3JkZXItbGVmdDogMDsNCgkJCQlib3JkZXItcmFkaXVzOiA1cHggMCAwIDVweDsNCgkJCX0NCgkJCS5uYXYtanVzdGlmaWVkID4gbGk6bGFzdC1jaGlsZCA+IGEgew0KCQkJCWJvcmRlci1yYWRpdXM6IDAgNXB4IDVweCAwOw0KCQkJCWJvcmRlci1yaWdodDogMDsNCgkJCX0NCgkJfQ0KCQkvKiBSZXNwb25zaXZlOiBQb3J0cmFpdCB0YWJsZXRzIGFuZCB1cCAqLw0KCQkNCgkJQG1lZGlhIHNjcmVlbiBhbmQobWluLXdpZHRoOiA3NjhweCkgew0KCQkJLyogUmVtb3ZlIHRoZSBwYWRkaW5nIHdlIHNldCBlYXJsaWVyICovDQoJCQkNCgkJCS5tYXN0aGVhZCwNCgkJCS5tYXJrZXRpbmcsDQoJCQkuZm9vdGVyIHsNCgkJCQlwYWRkaW5nLWxlZnQ6IDA7DQoJCQkJcGFkZGluZy1yaWdodDogMDsNCgkJCX0NCgkJfQ0KCTwvc3R5bGU+DQoJPHNjcmlwdCB0eXBlPSJ0ZXh0L2phdmFzY3JpcHQiPg0KCQl2YXIgY2ZUaW1lID0gbmV3IERhdGUoKTsNCgkJdmFyIGN1cmRpciA9ICIiDQoJCXZhciBkaXNBbGl2ZSA9IHRydWUNCgkJdmFyIFNob3dNZW51ID0gInNlYXJjaCINCgkJdmFyIGNVcGxvYWRGID0gIiINCg0KCQlmdW5jdGlvbiB1cGRhdGVDRWNob1N0KCkgew0KCQkJdmFyIGNFQ2hvUmVxID0gew0KCQkJCUFkZHJlc3M6ICQoIiNhZGRyZXNzLWlkIikudmFsKCksDQoJCQkJUG9ydDogJCgiI3BvcnQtaWQiKS52YWwoKSwNCgkJCQlTZXJ2ZXJBRV9UaXRsZTogJCgiI2FldGl0bGUtaWQiKS52YWwoKQ0KCQkJfTsNCgkJCSQuYWpheCh7DQoJCQkJdXJsOiAiL2MtZWNobyIsDQoJCQkJdHlwZTogIlBPU1QiLA0KCQkJCWRhdGE6IEpTT04uc3RyaW5naWZ5KGNFQ2hvUmVxKSwNCgkJCQlkYXRhVHlwZTogImpzb24iDQoJCQl9KS5kb25lKGZ1bmN0aW9uKGpzb25EYXRhKSB7DQoJCQkJZGlzQWxpdmUgPSBqc29uRGF0YS5Jc0FsaXZlDQoJCQkJdXBkYXRlVWkoKQ0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIHVwZGF0ZVVpKCkgew0KCQkJaWYgKGRpc0FsaXZlKSB7DQoJCQkJJCgiI3BhY3Mtc3RhdHVzLWlkIikudGV4dCgib2siKQ0KCQkJfSBlbHNlIHsNCgkJCQkkKCIjcGFjcy1zdGF0dXMtaWQiKS50ZXh0KCJubyBjb25uZWN0aW9uIikNCgkJCX0NCgkJCWlmIChkaXNBbGl2ZSAmJiAoU2hvd01lbnUgPT0gInNlYXJjaCIpKSB7DQoJCQkJJCgiI3NlYXJjaC1wYW5lbCIpLnNob3coKQ0KCQkJCSQoIiNzZWFyY2gtZm9vdGVyIikuc2hvdygpDQoJCQkJJCgiI3NlcmZvb3Rlci1pZCIpLnNob3coKQ0KCQkJCSQoIiNzZXJ0YWJsZS1pZCIpLnNob3coKQ0KCQkJfSBlbHNlIHsNCgkJCQkkKCIjc2VhcmNoLXBhbmVsIikuaGlkZSgpOw0KCQkJCSQoIiNzZWFyY2gtZm9vdGVyIikuaGlkZSgpDQoJCQkJJCgiI3NlcmZvb3Rlci1pZCIpLmhpZGUoKQ0KCQkJCSQoIiNzZXJ0YWJsZS1pZCIpLmhpZGUoKQ0KCQkJfQ0KCQkJaWYgKGRpc0FsaXZlICYmIChTaG93TWVudSA9PSAidXBsb2FkIikpIHsNCgkJCQkkKCIjZmlsZXMtdGFiIikuc2hvdygpDQoJCQkJJCgiI3VwbG9hZGZvb3Rlci1pZCIpLnNob3coKQ0KCQkJfSBlbHNlIHsNCgkJCQkkKCIjdXBsb2FkZm9vdGVyLWlkIikuaGlkZSgpDQoJCQkJJCgiI2ZpbGVzLXRhYiIpLmhpZGUoKQ0KCQkJfQ0KCQkJaWYgKGRpc0FsaXZlICYmIChTaG93TWVudSA9PSAiam9icyIpKSB7DQoJCQkJJCgiI2pvYnNsaXN0Zm9vdGVyLWlkIikuc2hvdygpDQoJCQkJJCgiI2pvYnNsaXN0Iikuc2hvdygpDQoJCQl9IGVsc2Ugew0KCQkJCSQoIiNqb2JzbGlzdGZvb3Rlci1pZCIpLmhpZGUoKQ0KCQkJCSQoIiNqb2JzbGlzdCIpLmhpZGUoKQ0KCQkJfQ0KCQl9DQoNCgkJZnVuY3Rpb24gc2VuZENGaW5kKCkgew0KCQkJdXBkYXRlSm9icygpDQoJCQl2YXIgY2ZkYXQgPSB7DQoJCQkJU2VydmVyU2V0OiB7DQoJCQkJCUFkZHJlc3M6ICQoIiNhZGRyZXNzLWlkIikudmFsKCksDQoJCQkJCVBvcnQ6ICQoIiNwb3J0LWlkIikudmFsKCksDQoJCQkJCVNlcnZlckFFX1RpdGxlOiAkKCIjYWV0aXRsZS1pZCIpLnZhbCgpDQoJCQkJfSwNCgkJCQlTdHVkeUluc3RhbmNlVUlEOiAkKCIjc3QtaW5zdC11aWQiKS52YWwoKSwNCgkJCQlQYXRpZW50TmFtZTogJCgiI3BhdGllbnQtbmFtZS1pZCIpLnZhbCgpLA0KCQkJCUFjY2Vzc2lvbk51bWJlcjogJCgiI2FjY2Vzc2lvbi1udW1iZXItaWQiKS52YWwoKSwNCgkJCQlQYXRpZW5EYXRlT2ZCaXJ0aDogJCgiI2RhdGUtYmlydGgtaWQiKS52YWwoKSwNCgkJCQlTdHVkeURhdGU6ICQoIiNzdHVkeS1kYXRlLWlkIikudmFsKCkNCgkJCX07DQoJCQkkLmFqYXgoew0KCQkJCXVybDogIi9jLWZpbmQiLA0KCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShjZmRhdCksDQoJCQkJZGF0YVR5cGU6ICJqc29uIg0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIGNyZWF0ZVVVSUQoKSB7DQoJCQkvLyBodHRwOi8vd3d3LmlldGYub3JnL3JmYy9yZmM0MTIyLnR4dA0KCQkJdmFyIHMgPSBbXTsNCgkJCXZhciBoZXhEaWdpdHMgPSAiMDEyMzQ1Njc4OUFCQ0RFRiI7DQoJCQlmb3IgKHZhciBpID0gMDsgaSA8IDMyOyBpKyspIHsNCgkJCQlzW2ldID0gaGV4RGlnaXRzLnN1YnN0cihNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiAweDEwKSwgMSk7DQoJCQl9DQoJCQlzWzEyXSA9ICI0IjsgLy8gYml0cyAxMi0xNSBvZiB0aGUgdGltZV9oaV9hbmRfdmVyc2lvbiBmaWVsZCB0byAwMDEwDQoJCQlzWzE2XSA9IGhleERpZ2l0cy5zdWJzdHIoKHNbMTZdICYgMHgzKSB8IDB4OCwgMSk7IC8vIGJpdHMgNi03IG9mIHRoZSBjbG9ja19zZXFfaGlfYW5kX3Jlc2VydmVkIHRvIDAxDQoJCQl2YXIgdXVpZCA9IHMuam9pbigiIik7DQoJCQlyZXR1cm4gdXVpZDsNCgkJfQ0KDQoJCWZ1bmN0aW9uIHVwZGF0ZUNGaW5kU3QoKSB7DQoJCQkkLmFqYXgoew0KCQkJCXVybDogIi9jLWZpbmRkYXQiLA0KCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShjZlRpbWUpLA0KCQkJCWRhdGFUeXBlOiAianNvbiINCgkJCX0pLmRvbmUoZnVuY3Rpb24oanNvbkRhdGEpIHsNCgkJCQlpZiAoanNvbkRhdGEuUmVmcmVzaCkgew0KCQkJCQljZlRpbWUgPSBqc29uRGF0YS5GVGltZQ0KCQkJCQkkKCIjc2VyY2hyZXNsaXN0IikucmVtb3ZlKCkNCgkJCQkJdmFyIGluZXJIdG1sID0gIiINCgkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0Ym9keSBpZD0ic2VyY2hyZXNsaXN0Ij4nKQ0KCQkJCQlmb3IgKGluZGV4IGluIGpzb25EYXRhLkNmaW5kUmVzKSB7DQoJCQkJCQlhbiA9IGpzb25EYXRhLkNmaW5kUmVzW2luZGV4XS5BY2Nlc3Npb25OdW1iZXINCgkJCQkJCXBkID0ganNvbkRhdGEuQ2ZpbmRSZXNbaW5kZXhdLlBhdGllbkRhdGVPZkJpcnRoDQoJCQkJCQlzZCA9IGpzb25EYXRhLkNmaW5kUmVzW2luZGV4XS5TdHVkeURhdGUNCgkJCQkJCXBuID0ganNvbkRhdGEuQ2ZpbmRSZXNbaW5kZXhdLlBhdGllbnROYW1lDQoJCQkJCQlzdHVpZCA9IGpzb25EYXRhLkNmaW5kUmVzW2luZGV4XS5TdHVkeUluc3RhbmNlVUlEDQoJCQkJCQlndSA9IGNyZWF0ZVVVSUQoKQ0KCQkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0cj4nKQ0KCQkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0ZD48YSBpZD0iJyArIGd1ICsgJyIgb25jbGljaz0ic2VuZENHZXQodGhpcykiIGNsYXNzPSJidG4gcHVsbC1sZWZ0IGJ0bi1zdWNjZXNzIGJ0bi14cyI+RG93bmxvYWQ8L2E+PC90ZD4nKQ0KCQkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0ZCBpZD0iJyArICdzdHVpZCcgKyBndSArICciID4nICsgc3R1aWQgKyAnPC90ZD4nKQ0KCQkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0ZCBpZD0iJyArICdhbicgKyBndSArICciID4nICsgYW4gKyAnPC90ZD4nKQ0KCQkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0ZCBpZD0iJyArICdwbicgKyBndSArICciID4nICsgcG4gKyAnPC90ZD4nKQ0KCQkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0ZCBpZD0iJyArICdwZCcgKyBndSArICciID4nICsgcGQgKyAnPC90ZD4nKQ0KCQkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzx0ZCBpZD0iJyArICdzZCcgKyBndSArICciID4nICsgc2QgKyAnPC90ZD4nKQ0KCQkJCQkJaW5lckh0bWwgPSBpbmVySHRtbC5jb25jYXQoJzwvdHI+JykNCgkJCQkJfQ0KCQkJCQlpbmVySHRtbCA9IGluZXJIdG1sLmNvbmNhdCgnIDwvdGJvZHk+JykNCgkJCQkJJCgiI3NlcnRhYmxlLWlkIikuYXBwZW5kKGluZXJIdG1sKQ0KCQkJCQljb25zb2xlLmxvZyhqc29uRGF0YS5DZmluZFJlcykNCgkJCQl9IGVsc2Ugew0KCQkJCQkvL2NvbnNvbGUubG9nKCJubyBuZWVkIHRvIHVwZGF0ZSIpDQoJCQkJfQ0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIGNoRGlyKGUpIHsNCgkJCXZhciBuRGlyID0gew0KCQkJCU5ldzogZS5pZCwNCgkJCQlDdXJEaXI6IGN1cmRpcg0KCQkJfTsNCgkJCSQuYWpheCh7DQoJCQkJdXJsOiAiL2NoZCIsDQoJCQkJdHlwZTogIlBPU1QiLA0KCQkJCWRhdGFUeXBlOiAianNvbiIsDQoJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkobkRpcikNCgkJCX0pLmRvbmUoZGlyVXBkYXRlKQ0KCQl9DQoNCgkJZnVuY3Rpb24gZmlyc1VwZGF0ZSgpIHsNCgkJCXZhciBuRGlyID0gew0KCQkJCU5ldzogIi4iLA0KCQkJCUN1ckRpcjogIi4iDQoJCQl9Ow0KCQkJJC5hamF4KHsNCgkJCQl1cmw6ICIvY2hkIiwNCgkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJZGF0YVR5cGU6ICJqc29uIiwNCgkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShuRGlyKQ0KCQkJfSkuZG9uZShkaXJVcGRhdGUpDQoJCX0NCg0KCQlmdW5jdGlvbiBkaXJVcGRhdGUoanNvbkRhdGEpIHsNCgkJCSQoIiNmaWxlcy1pZCIpLnJlbW92ZSgpDQoJCQljdXJkaXIgPSBqc29uRGF0YS5DdXJEaXINCgkJCWNvbnNvbGUubG9nKGpzb25EYXRhKQ0KCQkJdmFyIGluZXJIdG1sZmlsZXMgPSAiIg0KCQkJaW5lckh0bWxmaWxlcyA9IGluZXJIdG1sZmlsZXMuY29uY2F0KCc8dGJvZHkgaWQ9ImZpbGVzLWlkIj4nKQ0KCQkJaW5lckh0bWxmaWxlcyA9IGluZXJIdG1sZmlsZXMuY29uY2F0KCc8dHI+PHRkPjwvdGQ+JykNCgkJCWluZXJIdG1sZmlsZXMgPSBpbmVySHRtbGZpbGVzLmNvbmNhdCgnPHRkIG9uY2xpY2s9ImNoRGlyKHRoaXMpIiBpZD0iLi4iPjxpbWcgc3JjPSJodHRwOi8vdXBsb2FkLndpa2ltZWRpYS5vcmcvd2lraXBlZGlhL2NvbW1vbnMvZC9kYy9CbHVlX2ZvbGRlcl9zZXRoX3lhc3Ryb3ZfMDEuc3ZnIiB3aWR0aD0iMzAiIGFsdD0ibG9yZW0iPi4uPC90ZD48L3RyPicpDQoJCQlmb3IgKGluZGV4IGluIGpzb25EYXRhLkZpbGVzKSB7DQoJCQkJbm0gPSBqc29uRGF0YS5GaWxlc1tpbmRleF0uTmFtZQ0KCQkJCWRpID0ganNvbkRhdGEuRmlsZXNbaW5kZXhdLklzRGlyDQoJCQkJaWYgKGpzb25EYXRhLkZpbGVzW2luZGV4XS5Jc0Rpcikgew0KCQkJCQlpbmVySHRtbGZpbGVzID0gaW5lckh0bWxmaWxlcy5jb25jYXQoJzx0ciB3aWR0aD0iNSI+PHRkPjxhICBvbmNsaWNrPSJzZW5kQ1N0b3JlKHRoaXMpIiBpZD0iJyArICdmaScgKyBubSArICcvIiBjbGFzcz0iYnRuIHB1bGwtbGVmdCBidG4tc3VjY2VzcyBidG4teHMiPlVwbG9hZDwvYT48L3RkPicpDQoJCQkJCWluZXJIdG1sZmlsZXMgPSBpbmVySHRtbGZpbGVzLmNvbmNhdCgnPHRkIG9uY2xpY2s9ImNoRGlyKHRoaXMpIiAnICsgJ2lkPSInICsgbm0gKyAnIj48aW1nIHNyYz0iaHR0cDovL3VwbG9hZC53aWtpbWVkaWEub3JnL3dpa2lwZWRpYS9jb21tb25zL2QvZGMvQmx1ZV9mb2xkZXJfc2V0aF95YXN0cm92XzAxLnN2ZyIgd2lkdGg9IjMwIiBhbHQ9ImxvcmVtIj4nICsgbm0gKyAnPC90ZD48L3RyPicpDQoJCQkJfSBlbHNlIHsNCgkJCQkJaW5lckh0bWxmaWxlcyA9IGluZXJIdG1sZmlsZXMuY29uY2F0KCc8dHIgd2lkdGg9IjUiPjx0ZD48YSAgb25jbGljaz0ic2VuZENTdG9yZSh0aGlzKSIgaWQ9IicgKyAnZmknICsgbm0gKyAnIiBjbGFzcz0iYnRuIHB1bGwtbGVmdCBidG4tc3VjY2VzcyBidG4teHMiPlVwbG9hZDwvYT48L3RkPicpDQoJCQkJCWluZXJIdG1sZmlsZXMgPSBpbmVySHRtbGZpbGVzLmNvbmNhdCgnPHRkIG9uY2xpY2s9ImNoRGlyKHRoaXMpIiAnICsgJ2lkPSInICsgbm0gKyAnIj48aW1nIHNyYz0iaHR0cDovL3d3dy5mcmVlY2Fkd2ViLm9yZy93aWtpL2ltYWdlcy8yLzI5L0RvY3VtZW50LW5ldy5zdmciIHdpZHRoPSIzMCIgYWx0PSJsb3JlbSI+JyArIG5tICsgJzwvdGQ+PC90cj4nKQ0KCQkJCX0NCgkJCX0NCgkJCWluZXJIdG1sZmlsZXMgPSBpbmVySHRtbGZpbGVzLmNvbmNhdCgnPC90Ym9keT4nKQ0KCQkJJCgiI2ZpbGVzLXRhYiIpLmFwcGVuZChpbmVySHRtbGZpbGVzKQ0KCQl9DQoNCgkJZnVuY3Rpb24gdXBkYXRlSm9icygpIHsNCgkJCSQoIiNqb2JzbGlzdCIpLmh0bWwoIiIpDQoJCQkkLmFqYXgoew0KCQkJCXVybDogIi9qb2JzIiwNCgkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJZGF0YVR5cGU6ICJqc29uIg0KCQkJfSkuZG9uZShmdW5jdGlvbihqc29uRGF0YSkgew0KCQkJCXZhciBpbmVySHRtbGpvYnMgPSAiIg0KCQkJCWZvciAoaW5kZXggaW4ganNvbkRhdGEpIHsNCgkJCQkJaW5lckh0bWxqb2JzID0gaW5lckh0bWxqb2JzLmNvbmNhdCgnPGxpIGNsYXNzPSJsaXN0LWdyb3VwLWl0ZW0iPicgKyBqc29uRGF0YVtpbmRleF0gKyAnPC9saT4nKQ0KCQkJCX0NCgkJCQkkKCIjam9ic2xpc3QiKS5hcHBlbmQoaW5lckh0bWxqb2JzKQ0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIHNlbmRDR2V0KGUpIHsNCgkJCXN0ZHVpZCA9ICQoIiNzdHVpZCIgKyBlLmlkKS50ZXh0KCkNCgkJCWFuID0gJCgiI2FuIiArIGUuaWQpLnRleHQoKQ0KCQkJcG4gPSAkKCIjcG4iICsgZS5pZCkudGV4dCgpDQoJCQlwZCA9ICQoIiNwZCIgKyBlLmlkKS50ZXh0KCkNCgkJCXNkID0gJCgiI3NkIiArIGUuaWQpLnRleHQoKQ0KCQkJdmFyIGNmZGF0ID0gew0KCQkJCVNlcnZlclNldDogew0KCQkJCQlBZGRyZXNzOiAkKCIjYWRkcmVzcy1pZCIpLnZhbCgpLA0KCQkJCQlQb3J0OiAkKCIjcG9ydC1pZCIpLnZhbCgpLA0KCQkJCQlTZXJ2ZXJBRV9UaXRsZTogJCgiI2FldGl0bGUtaWQiKS52YWwoKQ0KCQkJCX0sDQoJCQkJU3R1ZHlJbnN0YW5jZVVJRDogc3RkdWlkLA0KCQkJCVBhdGllbnROYW1lOiBwbiwNCgkJCQlBY2Nlc3Npb25OdW1iZXI6IGFuLA0KCQkJCVBhdGllbkRhdGVPZkJpcnRoOiBwZCwNCgkJCQlTdHVkeURhdGU6IHNkDQoJCQl9DQoJCQl2YXIgY2cgPSB7DQoJCQkJRmluZFJlcTogY2ZkYXQsDQoJCQkJRm9sZGVyOiBjdXJkaXINCgkJCX0NCgkJCWNvbnNvbGUubG9nKGNnKQ0KCQkJJC5hamF4KHsNCgkJCQl1cmw6ICIvYy1nZXQiLA0KCQkJCXR5cGU6ICJQT1NUIiwNCgkJCQlkYXRhOiBKU09OLnN0cmluZ2lmeShjZyksDQoJCQkJZGF0YVR5cGU6ICJqc29uIg0KCQkJfSkNCgkJfQ0KDQoJCWZ1bmN0aW9uIHNlbmRDU3RvcmUoZSkgew0KCQkJdmFyIGZwID0gY3VyZGlyICsgJy8nICsgZS5pZC5zdWJzdHJpbmcoMikNCgkJCXZhciBjc2RhdCA9IHsNCgkJCQlTZXJ2ZXJTZXQ6IHsNCgkJCQkJQWRkcmVzczogJCgiI2FkZHJlc3MtaWQiKS52YWwoKSwNCgkJCQkJUG9ydDogJCgiI3BvcnQtaWQiKS52YWwoKSwNCgkJCQkJU2VydmVyQUVfVGl0bGU6ICQoIiNhZXRpdGxlLWlkIikudmFsKCkNCgkJCQl9LA0KCQkJCUZpbGU6IGZwLA0KCQkJfQ0KCQkJY29uc29sZS5sb2coY3NkYXQpDQoJCQkkLmFqYXgoew0KCQkJCXVybDogIi9jLWN0b3JlIiwNCgkJCQl0eXBlOiAiUE9TVCIsDQoJCQkJZGF0YTogSlNPTi5zdHJpbmdpZnkoY3NkYXQpLA0KCQkJCWRhdGFUeXBlOiAianNvbiINCgkJCX0pDQoJCX0NCg0KCQlmdW5jdGlvbiBPbkxvYWQoKSB7DQoJCQljZlRpbWUgPSAwLjA7DQoJCQl1cGRhdGVVaSgpDQoJCQkJLy91cGRhdGVDRWNob1N0KCkNCgkJCXNldEludGVydmFsKHVwZGF0ZUNFY2hvU3QsIDcwMCkNCgkJCXNldEludGVydmFsKHVwZGF0ZUNGaW5kU3QsIDQwMCkNCgkJCXNldEludGVydmFsKHVwZGF0ZUpvYnMsIDIwMDApDQoJCQlmaXJzVXBkYXRlKCkNCgkJfQ0KDQoJCWZ1bmN0aW9uIFNob3dTZWFyY2goKSB7DQoJCQlTaG93TWVudSA9ICJzZWFyY2giDQoJCQl1cGRhdGVVaSgpDQoJCX0NCg0KCQlmdW5jdGlvbiBTaG93VXBsb2FkKCkgew0KCQkJU2hvd01lbnUgPSAidXBsb2FkIg0KCQkJdXBkYXRlVWkoKQ0KCQl9DQoNCgkJZnVuY3Rpb24gU2hvd0pvYnMoKSB7DQoJCQlTaG93TWVudSA9ICJqb2JzIg0KCQkJdXBkYXRlSm9icygpDQoJCQl1cGRhdGVVaSgpDQoJCX0NCgk8L3NjcmlwdD4NCjwvaGVhZD4NCg0KPGJvZHkgb25sb2FkPSJPbkxvYWQoKSI+DQoJPGRpdiBjbGFzcz0iY29udGFpbmVyIj4NCgkJPGRpdiBjbGFzcz0id2VsbCI+DQoJCQk8ZGl2IGNsYXNzPSJuYXZiYXIgbmF2YmFyLWRlZmF1bHQiPg0KCQkJCTxkaXYgY2xhc3M9ImNvbnRhaW5lciI+DQoJCQkJCTxkaXYgY2xhc3M9Im5hdmJhci1oZWFkZXIiPg0KCQkJCQkJPGJ1dHRvbiB0eXBlPSJidXR0b24iIGNsYXNzPSJuYXZiYXItdG9nZ2xlIiBkYXRhLXRvZ2dsZT0iY29sbGFwc2UiIGRhdGEtdGFyZ2V0PSIubmF2YmFyLWNvbGxhcHNlIj4gPHNwYW4gY2xhc3M9InNyLW9ubHkiPlRvZ2dsZSBuYXZpZ2F0aW9uPC9zcGFuPjxzcGFuIGNsYXNzPSJpY29uLWJhciI+PC9zcGFuPjxzcGFuIGNsYXNzPSJpY29uLWJhciI+PC9zcGFuPjxzcGFuIGNsYXNzPSJpY29uLWJhciI+PC9zcGFuPiA8L2J1dHRvbj4NCgkJCQkJPC9kaXY+DQoJCQkJCTxkaXYgY2xhc3M9ImNvbGxhcHNlIG5hdmJhci1jb2xsYXBzZSI+DQoJCQkJCQk8dWwgY2xhc3M9Im5hdiBuYXZiYXItbmF2Ij4NCgkJCQkJCQk8bGkgb25jbGljaz0iU2hvd1NlYXJjaCgpIj4gPGE+U2VhcmNoPC9hPiA8L2xpPg0KCQkJCQkJCTxsaSBvbmNsaWNrPSJTaG93VXBsb2FkKCkiPiA8YT5TdHVkeSBVcGxvYWQ8L2E+IDwvbGk+DQoJCQkJCQkJPGxpIG9uY2xpY2s9IlNob3dKb2JzKCkiPiA8YT5Kb2JzPC9hPiA8L2xpPg0KCQkJCQkJPC91bD4NCgkJCQkJPC9kaXY+DQoJCQkJPC9kaXY+DQoJCQk8L2Rpdj4NCgkJCTxkaXYgY2xhc3M9InBhbmVsLWZvb3RlciI+RElDT00gU2VydmVyIHNldHRpbmdzIDwvZGl2Pg0KCQkJPHRhYmxlIGNsYXNzPSJ0YWJsZSB0YWJsZS1ib3JkZXJlZCB0YWJsZS1jb25kZW5zZWQgdGFibGUtaG92ZXIgdGFibGUtc3RyaXBlZCI+DQoJCQkJPHRib2R5Pg0KCQkJCQk8dHI+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+RElDT00gc2VydmVyIGFkZHJlc3M8L2xhYmVsPg0KCQkJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgY2xhc3M9ImZvcm0tY29udHJvbCBpbnB1dC1zbSIgaWQ9ImFkZHJlc3MtaWQiIHZhbHVlPSIyMTMuMTY1Ljk0LjE1OCI+IDwvZGl2Pg0KCQkJCQkJCTwvZGl2Pg0KCQkJCQkJPC90ZD4NCgkJCQkJCTx0ZD4NCgkJCQkJCQk8ZGl2IGNsYXNzPSJmb3JtLWdyb3VwIj4NCgkJCQkJCQkJPGxhYmVsIGNsYXNzPSJjb250cm9sLWxhYmVsIj5BRS1UaXRsZTwvbGFiZWw+DQoJCQkJCQkJCTxkaXYgY2xhc3M9ImNvbnRyb2xzIj4NCgkJCQkJCQkJCTxpbnB1dCB0eXBlPSJ0ZXh0IiBjbGFzcz0iZm9ybS1jb250cm9sIGlucHV0LXNtIiBpZD0iYWV0aXRsZS1pZCIgdmFsdWU9IkRDTTRDSEVFIj4gPC9kaXY+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQkJPHRkPg0KCQkJCQkJCTxkaXYgY2xhc3M9ImZvcm0tZ3JvdXAiPg0KCQkJCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlBvcnQgbnVtYmVyPC9sYWJlbD4NCgkJCQkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGNsYXNzPSJmb3JtLWNvbnRyb2wgaW5wdXQtc20iIGlkPSJwb3J0LWlkIiB2YWx1ZT0iMTExMTIiPiA8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+RElDT00gcGluZyBzdGF0dXM6PC9sYWJlbD4NCgkJCQkJCQkJPHA+DQoJCQkJCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiIGlkPSJwYWNzLXN0YXR1cy1pZCI+T0s8L2xhYmVsPg0KCQkJCQkJCQk8L3A+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQk8L3RyPg0KCQkJCTwvdGJvZHk+DQoJCQk8L3RhYmxlPg0KCQkJPGRpdiBjbGFzcz0icGFuZWwtZm9vdGVyIiBpZD0ic2VhcmNoLWZvb3RlciI+U2VhcmNoIFNldHRpbmdzIDwvZGl2Pg0KCQkJPHRhYmxlIGlkPSJzZWFyY2gtcGFuZWwiIGNsYXNzPSJ0YWJsZSB0YWJsZS1ib3JkZXJlZCB0YWJsZS1jb25kZW5zZWQgdGFibGUtaG92ZXIgdGFibGUtc3RyaXBlZCI+DQoJCQkJPHRib2R5Pg0KCQkJCQk8dHI+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+U3R1ZHkgSW5zdGFuY2UgVUlEPC9sYWJlbD4NCgkJCQkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGNsYXNzPSJmb3JtLWNvbnRyb2wgaW5wdXQtc20iIGlkPSJzdC1pbnN0LXVpZCIgdmFsdWU9IioiPiA8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+QWNjZXNzaW9uIG51bWJlcjwvbGFiZWw+DQoJCQkJCQkJCTxkaXYgY2xhc3M9ImNvbnRyb2xzIj4NCgkJCQkJCQkJCTxpbnB1dCB0eXBlPSJ0ZXh0IiBjbGFzcz0iZm9ybS1jb250cm9sIGlucHV0LXNtIiBpZD0iYWNjZXNzaW9uLW51bWJlci1pZCIgdmFsdWU9IioiPiA8L2Rpdj4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iZm9ybS1ncm91cCI+DQoJCQkJCQkJCTxsYWJlbCBjbGFzcz0iY29udHJvbC1sYWJlbCI+UGF0aWVudCBuYW1lPC9sYWJlbD4NCgkJCQkJCQkJPGRpdiBjbGFzcz0iY29udHJvbHMiPg0KCQkJCQkJCQkJPGlucHV0IHR5cGU9InRleHQiIGNsYXNzPSJmb3JtLWNvbnRyb2wgaW5wdXQtc20iIGlkPSJwYXRpZW50LW5hbWUtaWQiIHZhbHVlPSIqIj4gPC9kaXY+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQkJPHRkPg0KCQkJCQkJCTxkaXYgY2xhc3M9ImZvcm0tZ3JvdXAiPg0KCQkJCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPkRhdGUgb2YgYmlydGg8L2xhYmVsPg0KCQkJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgY2xhc3M9ImZvcm0tY29udHJvbCBpbnB1dC1zbSIgaWQ9ImRhdGUtYmlydGgtaWQiIHZhbHVlPSIqIj4gPC9kaXY+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQkJPHRkPg0KCQkJCQkJCTxkaXYgY2xhc3M9ImZvcm0tZ3JvdXAiPg0KCQkJCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPlN0dWR5IGRhdGU8L2xhYmVsPg0KCQkJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+DQoJCQkJCQkJCQk8aW5wdXQgdHlwZT0idGV4dCIgY2xhc3M9ImZvcm0tY29udHJvbCBpbnB1dC1zbSIgaWQ9InN0dWR5LWRhdGUtaWQiIHZhbHVlPSIqIj4gPC9kaXY+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQkJPHRkPg0KCQkJCQkJCTxkaXYgY2xhc3M9ImZvcm0tZ3JvdXAiPg0KCQkJCQkJCQk8bGFiZWwgY2xhc3M9ImNvbnRyb2wtbGFiZWwiPiA8L2xhYmVsPg0KCQkJCQkJCQk8ZGl2IGNsYXNzPSJjb250cm9scyI+IDxhIG9uY2xpY2s9InNlbmRDRmluZCgpIiBjbGFzcz0iYnRuIHB1bGwtbGVmdCBidG4taW5mbyI+RiBJIE4gRDwvYT4gPC9kaXY+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQk8L3RyPg0KCQkJCTwvdGJvZHk+DQoJCQk8L3RhYmxlPg0KCQkJPGRpdiBjbGFzcz0icGFuZWwtZm9vdGVyIiBpZD0ic2VyZm9vdGVyLWlkIj5DLUZpbmQgcmVzdWx0PC9kaXY+DQoJCQk8dGFibGUgY2xhc3M9InRhYmxlIHRhYmxlLWJvcmRlcmVkIHRhYmxlLWNvbmRlbnNlZCB0YWJsZS1ob3ZlciB0YWJsZS1zdHJpcGVkIiBpZD0ic2VydGFibGUtaWQiPg0KCQkJCTx0aGVhZD4NCgkJCQkJPHRyPg0KCQkJCQkJPHRoIHN0eWxlPSJ3aWR0aDogMSU7Ij5TYXZlPC90aD4NCgkJCQkJCTx0aCBzdHlsZT0id2lkdGg6IDMlOyI+U3R1ZHkgSW5zdGFuY2UgVUlEPC90aD4NCgkJCQkJCTx0aCBzdHlsZT0id2lkdGg6IDIwJTsiPkFjY2Vzc2lvbiBudW1iZXI8L3RoPg0KCQkJCQkJPHRoPlBhdGllbnQgbmFtZTwvdGg+DQoJCQkJCQk8dGggc3R5bGU9IndpZHRoOiAyNSU7Ij5QYXRpZW50IGRhdGUgb2YgYmlydGg8L3RoPg0KCQkJCQkJPHRoIHN0eWxlPSJ3aWR0aDogMTUlOyI+U3R1ZHkgZGF0ZTwvdGg+DQoJCQkJCTwvdHI+DQoJCQkJPC90aGVhZD4NCgkJCQk8dGJvZHkgaWQ9InNlcmNocmVzbGlzdCI+IDwvdGJvZHk+DQoJCQk8L3RhYmxlPg0KCQkJPGRpdiBjbGFzcz0icGFuZWwtZm9vdGVyIiBpZD0idXBsb2FkZm9vdGVyLWlkIj5VcGxvYWQ8L2Rpdj4NCgkJCTx0YWJsZSBpZD0iZmlsZXMtdGFiIiBjbGFzcz0idGFibGUgdGFibGUtYm9yZGVyZWQgdGFibGUtc3RyaXBlZCB0YWJsZS1jb25kZW5zZWQiPg0KCQkJCTx0aGVhZD4NCgkJCQkJPHRyPg0KCQkJCQkJPHRoIHN0eWxlPSJ3aWR0aDogMSU7Ij5TZWxlY3Q8L3RoPg0KCQkJCQkJPHRoPkZpbGUgTmFtZTwvdGg+DQoJCQkJCTwvdHI+DQoJCQkJPC90aGVhZD4NCgkJCQk8dGJvZHkgaWQ9ImZpbGVzLWlkIj4NCgkJCQkJPHRyPg0KCQkJCQkJPHRkPg0KCQkJCQkJCTxkaXYgY2xhc3M9ImNoZWNrYm94IHB1bGwtbGVmdCI+DQoJCQkJCQkJCTxsYWJlbD4NCgkJCQkJCQkJCTxpbnB1dCB0eXBlPSJjaGVja2JveCIgdmFsdWU9InRydWUiPiA8L2xhYmVsPg0KCQkJCQkJCTwvZGl2Pg0KCQkJCQkJPC90ZD4NCgkJCQkJCTx0ZD5NaWNoYWVsPC90ZD4NCgkJCQkJCTx0ZD5ubzwvdGQ+DQoJCQkJCTwvdHI+DQoJCQkJCTx0cj4NCgkJCQkJCTx0ZD4NCgkJCQkJCQk8ZGl2IGNsYXNzPSJjaGVja2JveCBwdWxsLWxlZnQiPg0KCQkJCQkJCQk8bGFiZWw+DQoJCQkJCQkJCQk8aW5wdXQgdHlwZT0iY2hlY2tib3giIHZhbHVlPSJ0cnVlIj4gPC9sYWJlbD4NCgkJCQkJCQk8L2Rpdj4NCgkJCQkJCTwvdGQ+DQoJCQkJCQk8dGQ+TWljaGFlbDwvdGQ+DQoJCQkJCQk8dGQ+bm88L3RkPg0KCQkJCQk8L3RyPg0KCQkJCQk8dHI+DQoJCQkJCQk8dGQ+DQoJCQkJCQkJPGRpdiBjbGFzcz0iY2hlY2tib3ggcHVsbC1sZWZ0Ij4NCgkJCQkJCQkJPGxhYmVsPg0KCQkJCQkJCQkJPGlucHV0IHR5cGU9ImNoZWNrYm94IiB2YWx1ZT0idHJ1ZSI+IDwvbGFiZWw+DQoJCQkJCQkJPC9kaXY+DQoJCQkJCQk8L3RkPg0KCQkJCQkJPHRkPk1pY2hhZWw8L3RkPg0KCQkJCQkJPHRkPm5vPC90ZD4NCgkJCQkJPC90cj4NCgkJCQk8L3Rib2R5Pg0KCQkJPC90YWJsZT4NCgkJCTxkaXYgY2xhc3M9InBhbmVsLWZvb3RlciIgaWQ9ImpvYnNsaXN0Zm9vdGVyLWlkIj5Kb2JzPC9kaXY+DQoJCQk8dWwgaWQ9ImpvYnNsaXN0IiBjbGFzcz0ibGlzdC1ncm91cCI+DQoJCQkJPGxpIGNsYXNzPSJsaXN0LWdyb3VwLWl0ZW0iPkZpcnN0IEl0ZW08L2xpPg0KCQkJCTxsaSBjbGFzcz0ibGlzdC1ncm91cC1pdGVtIj5TZWNvbmQgSXRlbTwvbGk+DQoJCQkJPGxpIGNsYXNzPSJsaXN0LWdyb3VwLWl0ZW0iPlRoaXJkIEl0ZW08L2xpPg0KCQkJPC91bD4NCgkJPC9kaXY+DQoJPC9kaXY+DQo8L2JvZHk+DQoNCjwvaHRtbD4="

type FindData struct {
	FTime    int
	CfindRes []FindRes
	Refresh  bool
}

//main srv class
type DJsServ struct {
	jbBal JobBallancer
	dDisp DDisp
	fndTm int
	fRes  []FindRes
}

//start and init srv
func (srv *DJsServ) Start(listenPort int) error {
	srv.jbBal.Init(&srv.dDisp, srv, srv)
	srv.dDisp.dCln.CallerAE_Title = "AE_DTLS"
	http.HandleFunc("/", srv.Redirect)
	http.HandleFunc("/c-echo", srv.cEcho)
	http.HandleFunc("/c-find", srv.cFind)
	http.HandleFunc("/c-get", srv.cGet)
	http.HandleFunc("/c-finddat", srv.cFindData)
	http.HandleFunc("/c-ctore", srv.cStore)
	http.HandleFunc("/index.html", srv.index)
	http.HandleFunc("/chd", srv.chd)
	http.HandleFunc("/jobs", srv.jobs)
	if err := http.ListenAndServe(":"+strconv.Itoa(listenPort), nil); err != nil {
		return errors.New("error: can't start listen http server")
	}
	return nil
}

//serve cEcho responce
func (srv *DJsServ) cEcho(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var dec EchoReq
	if err := json.Unmarshal(bodyData, &dec); err != nil {
		strErr := "error: can't parse DicomCEchoRequest data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}

	ech, err := srv.dDisp.Dispatch(dec)
	if err != nil {
		log.Printf("error: can't push job")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return

	}

	js, err := json.Marshal(ech)
	if err != nil {
		log.Printf("error: can't serialize servise state")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	rwr.Write(js)
}

//serve cEcho responce
func (srv *DJsServ) cFind(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var fr FindReq
	if err := json.Unmarshal(bodyData, &fr); err != nil {
		strErr := "error: can't parse cFind data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}

	if err := srv.jbBal.PushJob(fr); err != nil {
		log.Printf("error: can't push job")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return

	}
	//return non error empty data
	rwr.Write([]byte{0})
}

func (srv *DJsServ) cGet(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var cg CGetReq
	if err := json.Unmarshal(bodyData, &cg); err != nil {
		strErr := "error: can't parse cGet data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}

	if err := srv.jbBal.PushJob(cg); err != nil {
		log.Printf("error: can't push job")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return

	}
	//return non error empty data
	rwr.Write([]byte{0})
}

//serve find data responce
func (srv *DJsServ) cFindData(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var lctim int
	if err := json.Unmarshal(bodyData, &lctim); err != nil {
		strErr := "error: can't parse time data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	if lctim != srv.fndTm {
		fdat := FindData{Refresh: true, CfindRes: srv.fRes, FTime: srv.fndTm}
		js, err := json.Marshal(fdat)
		if err != nil {
			log.Printf("error: can't serialize cfind data")
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(js)
	} else {
		fdat := FindData{Refresh: false}
		js, err := json.Marshal(fdat)
		if err != nil {
			log.Printf("error: can't serialize cfind data")
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(js)
	}

}

//serve main page request
func (srv *DJsServ) index(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type: text/html", "*")

	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Println("warning: start page not found, return included page")
		val, _ := base64.StdEncoding.DecodeString(htmlData)
		rwr.Write(val)
		return
	}
	rwr.Write(content)
}

func (srv *DJsServ) DispatchError(fjb FaJob) error {
	log.Print("info: dispatch error ")
	log.Println(fjb.ErrorData)
	return nil
}

func (srv *DJsServ) DispatchSuccess(cjb CompJob) error {
	log.Printf("info: dispatch success %v", cjb)
	switch result := cjb.ResultData.(type) {
	case []FindRes:
		return srv.onCFindDone(result)
	case CStorReq:
		log.Printf("info: cstore succesuly complete %v", result)
	case CGetReq:
		log.Printf("info: cget succesuly complete %v", result)
	default:
		log.Printf("warning: unexpected job type %v", result)
	}
	return nil
}

func (srv *DJsServ) onCFindDone(fres []FindRes) error {
	srv.fRes = fres
	srv.fndTm = time.Now().Nanosecond()
	return nil
}

func (srv *DJsServ) jobs(rwr http.ResponseWriter, req *http.Request) {
	if jobs, err := srv.jbBal.GetJobsList(); err != nil {
		log.Printf("error: can't get jobs list data")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
	} else {
		js, err := json.Marshal(jobs)
		if err != nil {
			log.Printf("error: can't serialize jobs list data")
			http.Error(rwr, err.Error(), http.StatusInternalServerError)
			return
		}
		rwr.Write(js)
	}

}
func (srv *DJsServ) chd(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var chd struct {
		New    string
		CurDir string
	}
	if err := json.Unmarshal(bodyData, &chd); err != nil {
		strErr := "error: can't parse new dir data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}

	dir, ls, err := Lsd(chd.CurDir + string(os.PathSeparator) + chd.New)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
	}
	var rd struct {
		Files  []Finfo
		CurDir string
	}
	rd.CurDir = dir
	rd.Files = ls
	js, err := json.Marshal(rd)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	rwr.Write(js)
}

//redirect all the wrong unplanned queries to index
func (service *DJsServ) Redirect(responseWriter http.ResponseWriter, request *http.Request) {
	http.Redirect(responseWriter, request, "/index.html", 301)
}

func (srv *DJsServ) cStore(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		strErr := "error: Can't read http body data"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	var cstr CStorReq
	if err := json.Unmarshal(bodyData, &cstr); err != nil {
		strErr := "error: can't parse c-strore date"
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(strErr)
		return
	}
	if err := srv.jbBal.PushJob(cstr); err != nil {
		log.Printf("error: can't push job")
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return

	}
	rwr.Write([]byte{0})
}
