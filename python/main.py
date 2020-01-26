# import the necessary packages
from vision import Vision
from controller import Controller
from game import Game


def main():
    """
    Main method of the program.
    """

    vision = Vision()
    controller = Controller()
    game = Game(vision, controller)

    game.run()


if __name__ == '__main__':
    main()
