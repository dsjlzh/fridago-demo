'use strict';

Interceptor.attach(Module.getExportByName(null, 'open'), {
  onEnter(args) {
    console.log('[*] open(\"' + args[0].readUtf8String() + '\")');
  }
});
Interceptor.attach(Module.getExportByName(null, 'close'), {
  onEnter(args) {
    console.log('[*] close(' + args[0].toInt32() + ')');
  }
});
