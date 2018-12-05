$NetBSD: patch-sdext_source_pdfimport_xpdfwrapper_pdfioutdev__gpl.hxx,v 1.3 2018/12/04 12:56:49 ryoon Exp $

--- sdext/source/pdfimport/xpdfwrapper/pdfioutdev_gpl.hxx.orig	2018-10-29 19:55:29.000000000 +0000
+++ sdext/source/pdfimport/xpdfwrapper/pdfioutdev_gpl.hxx
@@ -151,22 +151,22 @@ namespace pdfi
 
         // Does this device use upside-down coordinates?
         // (Upside-down means (0,0) is the top left corner of the page.)
-        virtual GBool upsideDown() override { return gTrue; }
+        virtual bool upsideDown() override { return true; }
 
         // Does this device use drawChar() or drawString()?
-        virtual GBool useDrawChar() override { return gTrue; }
+        virtual bool useDrawChar() override { return true; }
 
         // Does this device use beginType3Char/endType3Char?  Otherwise,
         // text in Type 3 fonts will be drawn with drawChar/drawString.
-        virtual GBool interpretType3Chars() override { return gFalse; }
+        virtual bool interpretType3Chars() override { return false; }
 
         // Does this device need non-text content?
-        virtual GBool needNonText() override { return gTrue; }
+        virtual bool needNonText() override { return true; }
 
         //----- initialization and control
 
         // Set default transform matrix.
-        virtual void setDefaultCTM(double *ctm) override;
+        virtual void setDefaultCTM(const double *ctm) override;
 
         // Start a page.
         virtual void startPage(int pageNum, GfxState *state
@@ -233,40 +233,40 @@ namespace pdfi
 
         //----- image drawing
         virtual void drawImageMask(GfxState *state, Object *ref, Stream *str,
-                                   int width, int height, GBool invert,
+                                   int width, int height, bool invert,
 #if POPPLER_CHECK_VERSION(0, 12, 0)
-                                   GBool interpolate,
+                                   bool interpolate,
 #endif
-                                   GBool inlineImg) override;
+                                   bool inlineImg) override;
         virtual void drawImage(GfxState *state, Object *ref, Stream *str,
                                int width, int height, GfxImageColorMap *colorMap,
 #if POPPLER_CHECK_VERSION(0, 12, 0)
-                               GBool interpolate,
+                               bool interpolate,
 #endif
-                               int *maskColors, GBool inlineImg) override;
+                               int *maskColors, bool inlineImg) override;
         virtual void drawMaskedImage(GfxState *state, Object *ref, Stream *str,
                                      int width, int height,
                                      GfxImageColorMap *colorMap,
 #if POPPLER_CHECK_VERSION(0, 12, 0)
-                                     GBool interpolate,
+                                     bool interpolate,
 #endif
                                      Stream *maskStr, int maskWidth, int maskHeight,
-                                     GBool maskInvert
+                                     bool maskInvert
 #if POPPLER_CHECK_VERSION(0, 12, 0)
-                                     , GBool maskInterpolate
+                                     , bool maskInterpolate
 #endif
                                     ) override;
         virtual void drawSoftMaskedImage(GfxState *state, Object *ref, Stream *str,
                                          int width, int height,
                                          GfxImageColorMap *colorMap,
 #if POPPLER_CHECK_VERSION(0, 12, 0)
-                                         GBool interpolate,
+                                         bool interpolate,
 #endif
                                          Stream *maskStr,
                                          int maskWidth, int maskHeight,
                                          GfxImageColorMap *maskColorMap
 #if POPPLER_CHECK_VERSION(0, 12, 0)
-                                         , GBool maskInterpolate
+                                         , bool maskInterpolate
 #endif
                                         ) override;
 
