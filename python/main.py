# import the necessary packages
import argparse
import imutils
import numpy as np
import cv2  # Import the OpenCV library
import sys  # Enables the passing of arguments


def main():
    """
    Main method of the program.
    """

    # construct the argument parser and parse the arguments
    ap = argparse.ArgumentParser()
    ap.add_argument("-i", "--image", required=True, help="path to input image")
    args = vars(ap.parse_args())

    # load the input image (whose path was supplied via command line
    # argument) and display the image to our screen
    image = cv2.imread(args["image"])
    if image is None:
        print('Failed to load image file:', args["image"])
        sys.exit(1)

    # loop over the rotation angles
    for angle in np.arange(0, 360, 15):
        rotated = imutils.rotate(image, angle)
        cv2.imshow("Rotated (Problematic)", rotated)
        cv2.waitKey(0)

    # loop over the rotation angles again, this time ensuring
    # no part of the image is cut off
    for angle in np.arange(0, 360, 15):
        rotated = imutils.rotate_bound(image, angle)
        cv2.imshow("Rotated (Correct)", rotated)
        cv2.waitKey(0)


if __name__ == '__main__':
    main()
    cv2.destroyAllWindows()
