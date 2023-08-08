Task-2: Convert the image intro grayscale, display the grayscale image and save it
Code:
grayCBap = color.rgb2gray(CBap)
io.imshow(grayCBap)
io.imsave('Gray_CB_ap.jpg', grayCBap)


Task-3: Apply a 3x3 Median Filter on the image, display the filtered image and save it
Code:
import cv2
MedFilt = cv2.medianBlur(CBap, 3)
io.imshow(MedFilt)
io.imsave('Medfilt_CB_ap.jpg', MedFilt)


Task-4: Apply Histogram equalization on the image, display the image and save it
Code:
img =cv2.imread('/content/l-intro-1672278324.jpg')
IMG = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
src_eqlzd = cv2.equalizeHist(IMG)
io.imshow(src_eqlzd)
io.imsave('Eq.jpg',src_eqlzd)