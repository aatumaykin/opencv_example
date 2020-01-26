# import the necessary packages
import argparse
import cv2  # Import the OpenCV library
import imutils

minArea = 400


def main():
    """
    Main method of the program.
    """

    # construct the argument parser and parse the arguments
    ap = argparse.ArgumentParser()
    ap.add_argument("-v", "--video", required=True, help="Path to the video")
    args = vars(ap.parse_args())

    # load the image and convert it to grayscale
    cap = cv2.VideoCapture(args["video"])

    kernel = cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (3, 3))
    fgbg = cv2.createBackgroundSubtractorMOG2()

    while True:
        # Get webcam images
        ret, frame = cap.read()

        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        blurred = cv2.GaussianBlur(gray, (3, 3), 0)

        fgmask = fgbg.apply(blurred)
        fgmask = cv2.morphologyEx(fgmask, cv2.MORPH_OPEN, kernel)

        # threshold the image by setting all pixel values less than 225
        # to 255 (white; foreground) and all pixel values >= 225 to 255
        # (black; background), thereby segmenting the image
        thresh = cv2.threshold(fgmask, 30, 255, cv2.THRESH_BINARY)[1]

        # similarly, dilations can increase the size of the ground objects
        dilate = cv2.dilate(thresh, None, iterations=2)

        # find contours (i.e., outlines) of the foreground objects in the
        # thresholded image
        contours = cv2.findContours(dilate, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
        contours = imutils.grab_contours(contours)

        # loop over the contours
        for contour in contours:
            area = cv2.contourArea(contour)
            if area < minArea:
                continue

            (x, y, w, h) = cv2.boundingRect(contour)
            cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)

        cv2.imshow('Frame', frame)

        k = cv2.waitKey(30)
        if k == 27 or k == ord('q'):
            break

    cap.release()


if __name__ == '__main__':
    main()
    cv2.destroyAllWindows()
