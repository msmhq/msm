$(function() {
  var milestonesUrl = "https://api.github.com/repos/msmhq/msm/milestones?callback=?",
      tagsUrl       = "https://api.github.com/repos/msmhq/msm/tags?callback=?"
      stargazersUrl = "https://api.github.com/repos/msmhq/msm/stargazers?per_page=100&callback=?",
      contributorsUrl = "https://api.github.com/repos/msmhq/msm/contributors?per_page=100&callback=?";

  var maintainers = [
    "marcuswhybrow",
  ];
  var donators = [
    "joecabezas",
  ];

  function populateContributors() {
    var contributors = getAllPages(contributorsUrl);

    // Sort contributors by the number of contributins (most first)
    contributors.sort(function(a, b) {
      var aCount = a.contributions;
      var bCount = b.contributions;
      return (aCount < bCount) ? 1 : (aCount > bCount) ? -1 : 0;
    });

    $.each(contributors, function(i, contributor) {
      var $link = $("<a>")
        .attr("target", "_blank")
        .attr("href", "https://github.com/" + contributor.login + "/")
        .attr("ref", "tooltip")
        .attr("title", contributor.login)
        .addClass("contributor ga-track")
        .appendTo("#contributors")
        .tooltip({
          placement: "bottom",
        });
      $('<span>')
        .text('#' + contributor.contributions)
        .appendTo($link);
      $("<img>")
        .attr("src", contributor.avatar_url)
        .appendTo($link);

      if (maintainers.indexOf(contributor.login) >= 0) {
        $link.addClass("maintainer");
      }
      if (donators.indexOf(contributor.login) >= 0) {
        $link.addClass("donator");
      }
    });
  }

  function populateStargazers() {
    var stargazers = getAllPages(stargazersUrl);;
    $.each(stargazers, function(i, stargazer) {
      var $link = $("<a>")
        .attr("target", "_blank")
        .attr("href", "https://github.com/" + stargazer.login + "/")
        .attr("ref", "tooltip")
        .attr("title", stargazer.login)
        .addClass("stargazer ga-track")
        .appendTo("#stargazers")
        .tooltip({
          placement: "bottom",
        });
      $("<img>")
        .attr("src", stargazer.avatar_url)
        .appendTo($link);

      if (maintainers.indexOf(stargazer.login) >= 0) {
        $link.addClass("maintainer");
      }
      if (donators.indexOf(stargazer.login) >= 0) {
        $link.addClass("donator");
      }
    });
    $('.stargazers').slideDown();
  }

  function getAllPages(url, output, page) {
    if (typeof(output) === 'undefined') {
      output = [];
    }
    if (typeof(page) === 'undefined') {
      page = 1;
    }
    $.ajax({
      url: url + "&page=" + page,
      dataType: 'json',
      async: false,
      xhrFields: {
        withCredentials: true,
      },
      success: function(data) {
        if (data.length == 100) {
          output = getAllPages(url, output.concat(data), page + 1);
        } else {
          output = output.concat(data);
        }
      },
    });
    return output;
  }

  populateContributors();
  populateStargazers();

  $.getJSON(milestonesUrl, function(result) {
    var milestones = result.data;
    
    $.each(milestones, function(i, milestone) {
      var closedIssues = milestone.closed_issues;
      var totalIssues = closedIssues + milestone.open_issues;
      var milestonePercentage = closedIssues / totalIssues * 100;
      
      var $div = $("<div>").addClass("milestone");
      var $link = $("<a>")
            .attr("href", "https://github.com/msmhq/msm/issues?state=open&milestone=" + milestone.number)
            .addClass("ga-track")
            .appendTo($div);
      $("<h2>").text(milestone.title).appendTo($link);
      $("<p>").text(milestone.description).appendTo($div);
      var $progress = $("<div>").addClass("progress").appendTo($div);
      var $bar = $("<div>").addClass("bar").css("width", milestonePercentage + "%").appendTo($progress);
      $(".milestones").append($div);
    });
    
    $(".milestones").slideDown();
  });

  function compareTags(a, b) {
    if (a[0] > b[0]) return 1;
    else if (a[0] == b[0] && a[1] > b[1]) return 1;
    else if (a[0] == b[0] && a[1] == b[1] && a[2] > b[2]) return 1;
    else if (a[0] == b[0] && a[1] == b[1] && a[2] == b[2]) return 0;
    return -1;
  }

  function matchToTagVersion(match) {
    return [ parseInt(match[1], 10), parseInt(match[2], 10), parseInt(match[3], 10) ];
  }

  $.getJSON(tagsUrl, function(result) {
    var tags = result.data,
        latestTagVersion = [0, 0, 0],
        tagVersion, match;

    $.each(tags, function(i, tag) {
      if (! tag.name === 'undefined') {
        match = tag.name.match(/^(\d+)\.(\d+)\.(\d+)$/);
        if (match) {
          tagVersion = matchToTagVersion(match);
          if (compareTags(tagVersion, latestTagVersion) == 1)
            latestTagVersion = tagVersion;
        }
      }
    });

    var $installBtnVersion = $('#install-button small');
    match = $installBtnVersion.text().match(/^\((\d+)\.(\d+)\.(\d+)\)$/);
    if (match) {
      tagVersion = matchToTagVersion(match);
      if (compareTags(latestTagVersion, tagVersion) == 1)
        $installBtnVersion.text('(' + latestTagVersion.join('.') + ')');
    } 
  });

  $('.player').tooltip();
  $('.donator').tooltip({
    placement: 'bottom',
  });
});
