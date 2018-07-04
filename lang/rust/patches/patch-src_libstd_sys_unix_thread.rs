$NetBSD: patch-src_libstd_sys_unix_thread.rs,v 1.3 2018/06/25 13:31:11 ryoon Exp $

Fix stack-clash on illumos.

--- src/libstd/sys/unix/thread.rs.orig 2018-06-24 22:42:29.203295357 +0000
+++ src/libstd/sys/unix/thread.rs
@@ -309,7 +309,7 @@ pub mod guard {
 
         let stackaddr = get_stack_start_aligned()?;
 
-        if cfg!(target_os = "linux") {
+        if cfg!(any(target_os = "linux", target_os = "solaris")) {
             // Linux doesn't allocate the whole stack right away, and
             // the kernel has its own stack-guard mechanism to fault
             // when growing too close to an existing mapping.  If we map
@@ -326,13 +326,22 @@ pub mod guard {
             // Reallocate the last page of the stack.
             // This ensures SIGBUS will be raised on
             // stack overflow.
-            let result = mmap(stackaddr, PAGE_SIZE, PROT_NONE,
+            // Systems which enforce strict PAX MPROTECT do not allow
+            // to mprotect() a mapping with less restrictive permissions
+            // than the initial mmap() used, so we mmap() here with
+            // read/write permissions and only then mprotect() it to
+            // no permissions at all. See issue #50313.
+            let result = mmap(stackaddr, PAGE_SIZE, PROT_READ | PROT_WRITE,
                               MAP_PRIVATE | MAP_ANON | MAP_FIXED, -1, 0);
-
             if result != stackaddr || result == MAP_FAILED {
                 panic!("failed to allocate a guard page");
             }
 
+            let result = mprotect(stackaddr, PAGE_SIZE, PROT_NONE);
+            if result != 0 {
+                panic!("failed to protect the guard page");
+            }
+
             let guardaddr = stackaddr as usize;
             let offset = if cfg!(target_os = "freebsd") {
                 2
@@ -345,7 +354,7 @@ pub mod guard {
     }
 
     pub unsafe fn deinit() {
-        if !cfg!(target_os = "linux") {
+        if cfg!(not(any(target_os = "linux", target_os = "solaris"))) {
             if let Some(stackaddr) = get_stack_start_aligned() {
                 // Remove the protection on the guard page.
                 // FIXME: we cannot unmap the page, because when we mmap()
