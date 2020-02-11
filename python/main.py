# import the necessary packages
import argparse
import imutils
import cv2  # Import the OpenCV library
import sys  # Enables the passing of arguments
from logging import FileHandler
from vlogging import VisualRecord
import logging


def main():
    """
    Main method of the program.
    """

    # open the logging file
    logger = logging.getLogger("visual_logging_example")
    fh = FileHandler("demo.html", mode="w")
    # set the logger attributes
    logger.setLevel(logging.DEBUG)
    logger.addHandler(fh)

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

    logger.debug(VisualRecord("Image", image))

    # convert the image to grayscale
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    logger.debug(VisualRecord("Gray", gray))

    # applying edge detection we can find the outlines of objects in
    # images
    edged = cv2.Canny(gray, 30, 150)
    logger.debug(VisualRecord("Edged", edged))

    # threshold the image by setting all pixel values less than 225
    # to 255 (white; foreground) and all pixel values >= 225 to 255
    # (black; background), thereby segmenting the image
    thresh = cv2.threshold(gray, 225, 255, cv2.THRESH_BINARY_INV)[1]
    logger.debug(VisualRecord("Thresh", thresh))

    # find contours (i.e., outlines) of the foreground objects in the
    # thresholded image
    contours = cv2.findContours(thresh.copy(), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    contours = imutils.grab_contours(contours)
    output = image.copy()

    # loop over the contours
    for contour in contours:
        # draw each contour on the output image with a 3px thick purple
        # outline, then display the output contours one at a time
        cv2.drawContours(output, [contour], -1, (240, 0, 159), 3)
        logger.debug(VisualRecord("Contours", output))

    # draw the total number of contours found in purple
    text = "Found {} objects".format(len(contours))
    cv2.putText(output, text, (10, 25), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (240, 0, 159), 2)
    logger.debug(VisualRecord("Found objects", output))

    # we apply erosions to reduce the size of foreground objects
    erode = thresh.copy()
    erode = cv2.erode(erode, None, iterations=5)
    logger.debug(VisualRecord("Eroded", erode))

    # similarly, dilations can increase the size of the ground objects
    dilate = thresh.copy()
    dilate = cv2.dilate(dilate, None, iterations=5)
    logger.debug(VisualRecord("Dilated", dilate))

    # a typical operation we may want to apply is to take our mask and
    # apply a bitwise AND to our input image, keeping only the masked
    # regions
    mask = thresh.copy()
    output = cv2.bitwise_and(image, image, mask=mask)
    logger.debug(VisualRecord("Output", output))


if __name__ == '__main__':
    main()
