# import the necessary packages
import argparse
import imutils
import numpy as np
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

    # find all the 'black' shapes in the image
    lower = np.array([0, 0, 0])
    upper = np.array([15, 15, 15])
    shape_mask = cv2.inRange(image, lower, upper)

    # show the image

    cv2.namedWindow("Image", cv2.WINDOW_NORMAL)
    cv2.imshow("Image", image)

    # find the contours in the mask
    contours = cv2.findContours(shape_mask.copy(), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    contours = imutils.grab_contours(contours)

    print("I found {} black shapes".format(len(contours)))

    cv2.namedWindow("Mask", cv2.WINDOW_NORMAL)
    cv2.imshow("Mask", shape_mask)

    cv2.waitKey(0)

    # loop over the contours
    for c in contours:
        # draw the contour and show it
        cv2.drawContours(image, [c], -1, (255, 150, 150), 3)
        cv2.imshow("Image", image)

    cv2.waitKey(0)


if __name__ == '__main__':
    main()
    cv2.destroyAllWindows()
