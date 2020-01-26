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

    # define the list of boundaries
    boundaries = [
        ([15, 40, 140], [50, 90, 200]),
        ([130, 60, 15], [180, 80, 50]),
        ([0, 118, 130], [62, 150, 170]),
        ([40, 40, 20], [70, 50, 50])
    ]

    # load the input image (whose path was supplied via command line
    # argument) and display the image to our screen
    image = cv2.imread(args["image"])
    if image is None:
        print('Failed to load image file:', args["image"])
        sys.exit(1)

    cv2.namedWindow("Image", cv2.WINDOW_NORMAL)
    cv2.imshow("Image", image)

    cv2.waitKey(0)

    # loop over the boundaries
    for (lower, upper) in boundaries:
        # create NumPy arrays from the boundaries
        lower = np.array(lower, dtype="uint8")
        upper = np.array(upper, dtype="uint8")

        # find the colors within the specified boundaries and apply
        # the mask
        mask = cv2.inRange(image, lower, upper)
        output = cv2.bitwise_and(image, image, mask=mask)

        # show the images
        cv2.imshow("images", np.hstack([image, output]))
        cv2.waitKey(0)


if __name__ == '__main__':
    main()
    cv2.destroyAllWindows()
