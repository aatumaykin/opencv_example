# import the necessary packages
import argparse
import imutils
import cv2  # Import the OpenCV library


def main():
    """
    Main method of the program.
    """

    # construct the argument parser and parse the arguments
    ap = argparse.ArgumentParser()
    ap.add_argument("-i", "--image", required=True, help="path to input image")
    args = vars(ap.parse_args())

    # load the image, convert it to grayscale, blur it slightly,
    # and threshold it
    image = cv2.imread(args["image"])

    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    blurred = cv2.GaussianBlur(gray, (5, 5), 0)
    thresh = cv2.threshold(blurred, 60, 255, cv2.THRESH_BINARY)[1]

    # find contours in the thresholded image
    contours = cv2.findContours(thresh.copy(), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    contours = imutils.grab_contours(contours)

    # loop over the contours
    for c in contours:
        # compute the center of the contour
        m = cv2.moments(c)

        # draw the contour and center of the shape on the image
        cv2.drawContours(image, [c], -1, (0, 255, 0), 2)

        if m["m00"] > 0:
            c_x = int(m["m10"] / m["m00"])
            c_y = int(m["m01"] / m["m00"])
            cv2.circle(image, (c_x, c_y), 7, (255, 255, 255), -1)

    # show the image
    cv2.namedWindow("Image", cv2.WINDOW_NORMAL)
    cv2.imshow("Image", image)
    cv2.waitKey(0)


if __name__ == '__main__':
    main()
    cv2.destroyAllWindows()
